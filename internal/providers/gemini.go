package providers

import (
	"context"

	"google.golang.org/genai"
)

type GeminiProviderOptions struct {
	APIKey string
	Model  string
}

type GeminiProvider struct {
	client *genai.Client
	model  string
}

func CreateGeminiProvider(options GeminiProviderOptions) (*GeminiProvider, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  options.APIKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, err
	}

	return &GeminiProvider{client: client, model: options.Model}, nil
}

func (p *GeminiProvider) Prompt(ctx context.Context, promptOptions PromptOptions) (string, error) {
	response, err := p.client.Models.GenerateContent(ctx, p.model, genai.Text(promptOptions.UserMessage), &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(promptOptions.SystemMessage, genai.RoleUser),
	})
	if err != nil {
		return "", err
	}
	return response.Text(), nil
}
