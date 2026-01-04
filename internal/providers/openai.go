package providers

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIProviderOptions struct {
	APIKey string
	URL    string
	Model  string
}

type OpenAIProvider struct {
	client openai.Client
	model  string
}

func CreateOpenAIProvider(options OpenAIProviderOptions) (*OpenAIProvider, error) {
	client := openai.NewClient(
		option.WithAPIKey(options.APIKey),
		option.WithBaseURL(options.URL),
	)
	return &OpenAIProvider{client: client, model: options.Model}, nil
}

func (p *OpenAIProvider) Prompt(ctx context.Context, promptOptions PromptOptions) (string, error) {
	response, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(promptOptions.SystemMessage),
			openai.UserMessage(promptOptions.UserMessage),
		},
		Model: p.model,
	})
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}
