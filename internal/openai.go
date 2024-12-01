package internal

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient wraps the OpenAI client
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient initializes a new OpenAIClient
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
	}
}

// CreateChatCompletionRequest constructs the ChatCompletionRequest
func CreateChatCompletionRequest(prompt, diff string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: diff,
			},
		},
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		MaxTokens:        200,
		N:                1,
		Stream:           false,
	}
}

// SendChatCompletionRequest sends the request and returns the response
func (o *OpenAIClient) SendChatCompletionRequest(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return o.client.CreateChatCompletion(ctx, request)
}
