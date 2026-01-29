package providers

import (
	"context"
	"fmt"
	"sean111/althing/internal/tools"

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
	req := openrouter.ChatCompletionRequest{
		Model: p.model,
		Messages: []openrouter.ChatCompletionMessage{
			openrouter.UserMessage(promptOptions.UserMessage),
			openrouter.SystemMessage(promptOptions.SystemMessage),
		},
	}

	var functionDefs []openrouter.Tool
	for name, tool := range tools.ToolList {
		fd := openrouter.Tool{
			Type: openrouter.ToolTypeFunction,
			Function: &openrouter.FunctionDefinition{
				Name:        name,
				Description: tool.Description(),
				Parameters:  tool.Parameters(),
			},
		}
		functionDefs = append(functionDefs, fd)
	}

	req.Tools = functionDefs

	for i := 0; i < MAX_TOOL_CALLS; i++ {
		response, err := p.client.CreateChatCompletion(ctx, req)
		if err != nil {
			return "", err
		}

		message := response.Choices[0].Message

		if len(message.ToolCalls) == 0 {
			return message.Content.Text, nil
		}

		req.Messages = append(req.Messages, message)

		for _, toolCall := range message.ToolCalls {
			toolResponse, err := tools.ToolList[toolCall.Function.Name].Execute(ctx, toolCall.Function.Arguments)
			if err != nil {
				fmt.Printf("Error executing tool: %v\n", err)
				toolResponse = fmt.Sprintf("Error: %v", err)
			}
			req.Messages = append(req.Messages, openrouter.ToolMessage(toolResponse, toolCall.ID))
		}
	}
	return "", fmt.Errorf("reached maximum number of tool calls")
}
