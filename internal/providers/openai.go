package providers

import (
	"context"
	"fmt"
	"sean111/althing/internal/tools"

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
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(promptOptions.SystemMessage),
			openai.UserMessage(promptOptions.UserMessage),
		},
		Model: p.model,
	}

	var agentTools []openai.ChatCompletionToolUnionParam
	for name, tool := range tools.ToolList {
		t := openai.ChatCompletionFunctionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        name,
				Description: openai.String(tool.Description()),
				Parameters:  tool.Parameters(),
			},
		}
		agentTools = append(agentTools, openai.ChatCompletionToolUnionParam{OfFunction: &t})
	}

	params.Tools = agentTools

	response, err := p.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}

	toolCalls := response.Choices[0].Message.ToolCalls

	if len(toolCalls) == 0 {
		return response.Choices[0].Message.Content, nil
	}

	params.Messages = append(params.Messages, response.Choices[0].Message.ToParam())

	for _, toolCall := range toolCalls {
		toolResponse, err := tools.ToolList[toolCall.Function.Name].Execute(context.Background(), toolCall.Function.Arguments)
		if err != nil {
			fmt.Printf("Error executing tool: %v\n", err)
			continue
		}
		params.Messages = append(params.Messages, openai.ToolMessage(toolResponse, toolCall.ID))

	}

	response, err = p.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
