package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ycchen/ai-agent-go/internal/ollama"
)

type ReadFileTool struct{}

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Description() string {
	return "Reads the contents of a file from the filesystem"
}

func (t *ReadFileTool) Parameters() ollama.ParameterSchema {
	return ollama.ParameterSchema{
		Type: "object",
		Properties: map[string]ollama.PropertyDefinition{
			"file_path": {
				Type:        "string",
				Description: "Path to the file to read",
			},
		},
		Required: []string{"file_path"},
	}
}

func (t *ReadFileTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	filePath, ok := args["file_path"].(string)
	if !ok {
		return "", fmt.Errorf("file_path must be a string")
	}

	// Security: Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(filePath)

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

type WriteFileTool struct{}

func (t *WriteFileTool) Name() string {
	return "write_file"
}

func (t *WriteFileTool) Description() string {
	return "Writes content to a file, creating it if it doesn't exist, overwriting if it does"
}

func (t *WriteFileTool) Parameters() ollama.ParameterSchema {
	return ollama.ParameterSchema{
		Type: "object",
		Properties: map[string]ollama.PropertyDefinition{
			"file_path": {
				Type:        "string",
				Description: "Path to the file to write",
			},
			"content": {
				Type:        "string",
				Description: "Content to write to the file",
			},
		},
		Required: []string{"file_path", "content"},
	}
}

func (t *WriteFileTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	filePath, ok := args["file_path"].(string)
	if !ok {
		return "", fmt.Errorf("file_path must be a string")
	}

	content, ok := args["content"].(string)
	if !ok {
		return "", fmt.Errorf("content must be a string")
	}

	// Security: Clean the path
	cleanPath := filepath.Clean(filePath)

	// Create parent directories if they don't exist
	dir := filepath.Dir(cleanPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(cleanPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Successfully wrote %d bytes to %s", len(content), cleanPath), nil
}
