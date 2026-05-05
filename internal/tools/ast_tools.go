package tools

import (
	"context"
	"fmt"

	"github.com/ycchen/ai-agent-go/internal/ollama"
	astpkg "github.com/ycchen/ai-agent-go/internal/tools/ast"
)

type ExtractGoFunctionsTool struct {
	analyzer *astpkg.Analyzer
}

func NewExtractGoFunctionsTool() *ExtractGoFunctionsTool {
	return &ExtractGoFunctionsTool{
		analyzer: astpkg.NewAnalyzer(),
	}
}

func (t *ExtractGoFunctionsTool) Name() string {
	return "extract_go_functions"
}

func (t *ExtractGoFunctionsTool) Description() string {
	return "Extracts all function signatures from a Go source file, including parameters, return types, and documentation"
}

func (t *ExtractGoFunctionsTool) Parameters() ollama.ParameterSchema {
	return ollama.ParameterSchema{
		Type: "object",
		Properties: map[string]ollama.PropertyDefinition{
			"file_path": {
				Type:        "string",
				Description: "Path to the Go source file",
			},
		},
		Required: []string{"file_path"},
	}
}

func (t *ExtractGoFunctionsTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	filePath, ok := args["file_path"].(string)
	if !ok {
		return "", fmt.Errorf("file_path must be a string")
	}

	functions, err := t.analyzer.ExtractFunctions(filePath)
	if err != nil {
		return "", err
	}

	// Return formatted output
	return astpkg.FormatFunctions(functions), nil
}

type ExtractGoTypesTool struct {
	analyzer *astpkg.Analyzer
}

func NewExtractGoTypesTool() *ExtractGoTypesTool {
	return &ExtractGoTypesTool{
		analyzer: astpkg.NewAnalyzer(),
	}
}

func (t *ExtractGoTypesTool) Name() string {
	return "extract_go_types"
}

func (t *ExtractGoTypesTool) Description() string {
	return "Extracts all type definitions (structs, interfaces, type aliases) from a Go source file"
}

func (t *ExtractGoTypesTool) Parameters() ollama.ParameterSchema {
	return ollama.ParameterSchema{
		Type: "object",
		Properties: map[string]ollama.PropertyDefinition{
			"file_path": {
				Type:        "string",
				Description: "Path to the Go source file",
			},
		},
		Required: []string{"file_path"},
	}
}

func (t *ExtractGoTypesTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	filePath, ok := args["file_path"].(string)
	if !ok {
		return "", fmt.Errorf("file_path must be a string")
	}

	types, err := t.analyzer.ExtractTypes(filePath)
	if err != nil {
		return "", err
	}

	return astpkg.FormatTypes(types), nil
}

type AnalyzeGoPackageTool struct {
	analyzer *astpkg.Analyzer
}

func NewAnalyzeGoPackageTool() *AnalyzeGoPackageTool {
	return &AnalyzeGoPackageTool{
		analyzer: astpkg.NewAnalyzer(),
	}
}

func (t *AnalyzeGoPackageTool) Name() string {
	return "analyze_go_package"
}

func (t *AnalyzeGoPackageTool) Description() string {
	return "Analyzes an entire Go package directory and extracts all exported functions and types"
}

func (t *AnalyzeGoPackageTool) Parameters() ollama.ParameterSchema {
	return ollama.ParameterSchema{
		Type: "object",
		Properties: map[string]ollama.PropertyDefinition{
			"package_path": {
				Type:        "string",
				Description: "Path to the Go package directory",
			},
		},
		Required: []string{"package_path"},
	}
}

func (t *AnalyzeGoPackageTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	packagePath, ok := args["package_path"].(string)
	if !ok {
		return "", fmt.Errorf("package_path must be a string")
	}

	analysis, err := t.analyzer.AnalyzePackage(packagePath)
	if err != nil {
		return "", err
	}

	return astpkg.FormatPackageAnalysis(analysis), nil
}
