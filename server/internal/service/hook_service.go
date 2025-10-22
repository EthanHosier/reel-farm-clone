package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/ethanhosier/reel-farm/internal/api"
	"github.com/ethanhosier/reel-farm/internal/repository"
	"github.com/google/uuid"
)

const (
	promptTemplate = `
You are a helpful assistant that generates short hooks for a tiktok slideshow.
You will be given a prompt and you will need to generate a list of hooks for the prompt.

For example:
Prompt: "Plants dying in my house"
Hooks:
- "5 things I wish I knew before killing my plants"
- "um so why did it take a plant expert explaining to me that traditional planters are so expensive just to constantly water plants..."
- "fun fact you're probably spending way too much time and money on watering your plants when you don't need to"

The hooks should be returned in a json array of strings.
For example:
[
  "hook1",
  "hook2",
  "hook3"
]

The prompt is: {{.Prompt}}
The number of hooks to generate is: {{.NumHooks}}
`

	creditCost = 10
)

type HookService struct {
	userRepo   *repository.UserRepository
	hookRepo   *repository.HookRepository
	llmService *LLMService
}

type HookTemplateData struct {
	Prompt   string
	NumHooks int
}

type HookResponse struct {
	Hooks []string `json:"hooks"`
}

func NewHookService(userRepo *repository.UserRepository, hookRepo *repository.HookRepository, llmService *LLMService) *HookService {
	return &HookService{
		userRepo:   userRepo,
		hookRepo:   hookRepo,
		llmService: llmService,
	}
}

// TODO: Add idempotency and race condition protection
func (s *HookService) GenerateHooks(ctx context.Context, userID uuid.UUID, prompt string, numHooks int) ([]api.Hook, error) {
	// Use transaction to atomically check and deduct credits
	err := s.userRepo.WithTransaction(ctx, func(txRepo *repository.UserRepository) error {
		// Check if user has enough credits
		userAccount, err := txRepo.GetUserAccount(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user account: %w", err)
		}

		if userAccount.Credits < creditCost {
			return fmt.Errorf("insufficient credits: have %d, need %d", userAccount.Credits, creditCost)
		}

		// Remove credits
		err = txRepo.RemoveCreditsFromUser(ctx, userID, creditCost)
		if err != nil {
			return fmt.Errorf("failed to remove credits: %w", err)
		}

		return nil
	})
	if err != nil {
		_ = s.userRepo.AddCreditsToUser(ctx, userID, creditCost)
		return nil, fmt.Errorf("failed to deduct credits: %w", err)
	}

	hooks, err := s.doGenerateHooks(ctx, prompt, numHooks)
	if err != nil {
		_ = s.userRepo.AddCreditsToUser(ctx, userID, creditCost)
		return nil, fmt.Errorf("failed to generate hooks: %w", err)
	}

	// Store hooks in database and collect results
	generationID := uuid.New()
	createdHooks, err := s.hookRepo.CreateHooksBatch(ctx, userID, generationID, prompt, hooks, creditCost)
	if err != nil {
		// If storing fails, refund credits and return error
		_ = s.userRepo.AddCreditsToUser(ctx, userID, creditCost)
		return nil, fmt.Errorf("failed to store hooks: %w", err)
	}

	// Convert database hooks to API hooks
	var hookResults []api.Hook
	for _, dbHook := range createdHooks {
		hookResults = append(hookResults, api.Hook{
			Id:   dbHook.ID,
			Text: dbHook.HookText,
		})
	}

	return hookResults, nil
}

func (s *HookService) doGenerateHooks(ctx context.Context, prompt string, numHooks int) ([]string, error) {
	tmpl, err := template.New("hookPrompt").Parse(promptTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with the provided data
	var buf bytes.Buffer
	data := HookTemplateData{
		Prompt:   prompt,
		NumHooks: numHooks,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Get the generated prompt
	generatedPrompt := buf.String()

	// Call the LLM service
	response, err := s.llmService.GenerateText(ctx, generatedPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate text: %w", err)
	}

	// Parse the JSON response
	var hooks []string
	if err := json.Unmarshal([]byte(response), &hooks); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return hooks, nil
}

// DeleteHook deletes a hook (only if it belongs to the user)
func (s *HookService) DeleteHook(ctx context.Context, hookID uuid.UUID, userID uuid.UUID) error {
	err := s.hookRepo.DeleteHook(ctx, hookID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete hook: %w", err)
	}
	return nil
}
