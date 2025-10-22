package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"
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
	llmService *LLMService
	template   *template.Template
}

type HookTemplateData struct {
	Prompt   string
	NumHooks int
}

type HookResponse struct {
	Hooks []string `json:"hooks"`
}

func NewHookService(llmService *LLMService) *HookService {
	// Parse the template

	return &HookService{
		llmService: llmService,
	}
}

func (s *HookService) GenerateHooks(ctx context.Context, prompt string, numHooks int) ([]string, error) {
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
