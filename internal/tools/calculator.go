package tools

import (
	"context"
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/ycchen/ai-agent-go/internal/ollama"
)

type CalculatorTool struct{}

func (t *CalculatorTool) Name() string {
	return "calculator"
}

func (t *CalculatorTool) Description() string {
	return "Evaluates mathematical expressions safely. Supports +, -, *, /, (), etc."
}

func (t *CalculatorTool) Parameters() ollama.ParameterSchema {
	return ollama.ParameterSchema{
		Type: "object",
		Properties: map[string]ollama.PropertyDefinition{
			"expression": {
				Type:        "string",
				Description: "Mathematical expression to evaluate (e.g., '2 + 2', '(10 * 5) / 2')",
			},
		},
		Required: []string{"expression"},
	}
}

func (t *CalculatorTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	expression, ok := args["expression"].(string)
	if !ok {
		return "", fmt.Errorf("expression must be a string")
	}

	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", fmt.Errorf("invalid expression: %w", err)
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", fmt.Errorf("evaluation error: %w", err)
	}

	return fmt.Sprintf("%v", result), nil
}
