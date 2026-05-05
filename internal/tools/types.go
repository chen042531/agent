package tools

import (
	"context"

	"github.com/ycchen/ai-agent-go/internal/ollama"
)

type Tool interface {
	Name() string
	Description() string
	Parameters() ollama.ParameterSchema
	Execute(ctx context.Context, args map[string]interface{}) (string, error)
}
