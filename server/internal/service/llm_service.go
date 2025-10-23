package service

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

// LLMService handles LLM operations using OpenAI
type LLMService struct {
	client openai.Client
}

// NewLLMService creates a new LLM service
func NewLLMService() *LLMService {
	// Get OpenAI API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &LLMService{
		client: client,
	}
}

// GenerateText takes a prompt and returns the generated text response
func (s *LLMService) GenerateText(ctx context.Context, prompt string) (string, error) {
	// Create chat completion request
	chatCompletion, err := s.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String(prompt),
						},
					},
				},
			},
			Model: shared.ChatModelGPT5Mini,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate text: %w", err)
	}

	// Extract the response content
	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from OpenAI")
	}

	choice := chatCompletion.Choices[0]
	if choice.Message.Content == "" {
		return "", fmt.Errorf("no content in response message")
	}

	return choice.Message.Content, nil
}
