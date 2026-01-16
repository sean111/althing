package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sean111/althing/internal/formatting"

	"github.com/spf13/viper"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type Search struct {
}

func NewSearch() *Search {
	return &Search{}
}

func (s *Search) Name() string {
	return "web_search"
}

func (s *Search) Description() string {
	return "Search the web using Google"
}

func (s *Search) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "The search query to run",
			},
		},
		"required": []string{"query"},
	}
}

func (s *Search) Execute(ctx context.Context, jsonArgs string) (string, error) {
	var args struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal([]byte(jsonArgs), &args); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	apiKey := viper.GetString("tools.google-search.api_key")
	searchEngineId := viper.GetString("tools.google-search.search_engine_id")

	service, err := customsearch.NewService(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		return "", fmt.Errorf("failed to create custom search service: %w", err)
	}
	resp, err := service.Cse.List().Cx(searchEngineId).Q(args.Query).Do()

	if err != nil {
		return "", fmt.Errorf("failed to search: %w", err)
	}

	response := ""

	for _, item := range resp.Items {
		response += fmt.Sprintf("Title: %s\nLink: %s\nSnipplet:%s\n ", item.Title, item.Link, item.Snippet)
	}

	fmt.Printf("%s:\nQuery:\n%s\nResponse:\n%s\n", formatting.ToolCallType.Render("web_search"), formatting.ResponseStyle.Render(args.Query), formatting.ResponseStyle.Render(response))

	return response, nil

}
