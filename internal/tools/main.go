package tools

import "context"

var ToolList map[string]Tool

type Tool interface {
	Name() string
	Description() string
	Parameters() map[string]interface{}
	Execute(ctx context.Context, jsonArgs string) (string, error)
}
