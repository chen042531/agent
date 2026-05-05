package tools

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ycchen/ai-agent-go/internal/ollama"
)

type ExecuteCommandTool struct{}

func (t *ExecuteCommandTool) Name() string {
	return "execute_system_command"
}

func (t *ExecuteCommandTool) Description() string {
	return "Executes a system command and returns its output. Timeout: 30 seconds."
}

func (t *ExecuteCommandTool) Parameters() ollama.ParameterSchema {
	return ollama.ParameterSchema{
		Type: "object",
		Properties: map[string]ollama.PropertyDefinition{
			"command": {
				Type:        "string",
				Description: "The shell command to execute",
			},
		},
		Required: []string{"command"},
	}
}

func (t *ExecuteCommandTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	command, ok := args["command"].(string)
	if !ok {
		return "", fmt.Errorf("command must be a string")
	}

	// Create context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Execute command
	cmd := exec.CommandContext(cmdCtx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}
