package providers

import (
	"context"

	"github.com/revrost/go-openrouter"
)

type OpenRouterProviderOptions struct {
	Token string
	Title string
	URL   string
	Model string
}

type OpenRouterProvider struct {
	client *openrouter.Client
	model  string
}

func CreateOpenRouterProvider(options OpenRouterProviderOptions) (*OpenRouterProvider, error) {
	client := openrouter.NewClient(options.Token, openrouter.WithXTitle(options.Title), openrouter.WithHTTPReferer(options.URL))
	return &OpenRouterProvider{client: client, model: options.Model}, nil
}

func (p *OpenRouterProvider) Prompt(ctx context.Context, promptOptions PromptOptions) (string, error) {
	response, err := p.client.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model: p.model,
			Messages: []openrouter.ChatCompletionMessage{
				openrouter.UserMessage(promptOptions.UserMessage),
				openrouter.SystemMessage(promptOptions.SystemMessage),
			},
		},
	)

	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content.Text, nil
}
