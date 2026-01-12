package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sean111/althing/internal/formatting"

	"github.com/ajanicij/goduckgo/goduckgo"
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
	return "Search the web using DuckDuckGo"
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
	result, err := goduckgo.Query(args.Query)
	if err != nil {
		return "", fmt.Errorf("failed to search: %w", err)
	}

	var response string

	if result.AbstractText != "" {
		response += fmt.Sprintf("Instant Answer:\nAbstract: %s\nSource: %s\nURL: %s\n", result.AbstractText, result.AbstractSource, result.AbstractURL)

	}

	response += fmt.Sprintf("Related Topics:\n")

	for i, topic := range result.RelatedTopics {
		if i > 5 {
			break
		}
		response += fmt.Sprintf("Title: %s\nLink: %s\n", topic.Text, topic.FirstURL)
	}

	fmt.Printf("%s: Query:\n%s\n Response:\n%s\n", formatting.ToolCallType.Render("Search"), formatting.ResponseStyle.Render(args.Query), formatting.ResponseStyle.Render(response))

	return response, nil
}
