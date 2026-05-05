package main

import (
	"log"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Get workspace root from environment or use current directory
	workspaceRoot := os.Getenv("WORKSPACE_ROOT")
	if workspaceRoot == "" {
		var err error
		workspaceRoot, err = os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}
	}

	// Create MCP server
	s := server.NewMCPServer(
		"Go Code Analyzer MCP",
		"1.0.0",
	)

	analyzer := NewCodeAnalyzer(workspaceRoot)

	// Register all tools
	registerTools(s, analyzer)

	// Register resources
	registerResources(s, analyzer)

	// Start stdio server
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func registerTools(s *server.MCPServer, analyzer *CodeAnalyzer) {
	// Tool 1: List packages
	s.AddTool(
		mcp.NewTool("list_packages",
			mcp.WithDescription("List all Go packages in the workspace with optional summary level"),
			mcp.WithString("path",
				mcp.Description("Relative path from workspace root (default: current directory)"),
			),
			mcp.WithString("level",
				mcp.Description("Detail level: minimal (names only), summary (with type/func counts), detailed (with exports list)"),
				mcp.Enum("minimal", "summary", "detailed"),
			),
		),
		func(args map[string]interface{}) (*mcp.CallToolResult, error) {
			return analyzer.ListPackages(args)
		},
	)

	// Tool 2: Analyze package
	s.AddTool(
		mcp.NewTool("analyze_package",
			mcp.WithDescription("Analyze a specific Go package with configurable depth. Include parameter accepts array of: types, functions, imports, methods, fields"),
			mcp.WithString("package_path",
				mcp.Description("Path to the package directory"),
				mcp.Required(),
			),
			mcp.WithBoolean("exported_only",
				mcp.Description("Only include exported symbols (default: true)"),
			),
		),
		func(args map[string]interface{}) (*mcp.CallToolResult, error) {
			return analyzer.AnalyzePackage(args)
		},
	)

	// Tool 3: Get function details
	s.AddTool(
		mcp.NewTool("get_function_details",
			mcp.WithDescription("Get detailed information about a specific function or method"),
			mcp.WithString("file_path",
				mcp.Description("Path to the Go file"),
				mcp.Required(),
			),
			mcp.WithString("function_name",
				mcp.Description("Name of the function (use 'Type.Method' for methods)"),
				mcp.Required(),
			),
			mcp.WithBoolean("include_body",
				mcp.Description("Include function body source code"),
			),
		),
		func(args map[string]interface{}) (*mcp.CallToolResult, error) {
			return analyzer.GetFunctionDetails(args)
		},
	)

	// Tool 4: Get type details
	s.AddTool(
		mcp.NewTool("get_type_details",
			mcp.WithDescription("Get detailed information about a specific type (struct, interface, etc.)"),
			mcp.WithString("file_path",
				mcp.Description("Path to the Go file"),
				mcp.Required(),
			),
			mcp.WithString("type_name",
				mcp.Description("Name of the type"),
				mcp.Required(),
			),
			mcp.WithBoolean("include_methods",
				mcp.Description("Include methods of this type"),
			),
		),
		func(args map[string]interface{}) (*mcp.CallToolResult, error) {
			return analyzer.GetTypeDetails(args)
		},
	)

	// Tool 5: Search code
	s.AddTool(
		mcp.NewTool("search_code",
			mcp.WithDescription("Search for functions, types, or symbols in the codebase"),
			mcp.WithString("query",
				mcp.Description("Search query (name pattern)"),
				mcp.Required(),
			),
			mcp.WithString("kind",
				mcp.Description("Kind of symbol to search for"),
				mcp.Enum("all", "function", "type", "method", "interface"),
			),
			mcp.WithNumber("max_results",
				mcp.Description("Maximum number of results"),
			),
		),
		func(args map[string]interface{}) (*mcp.CallToolResult, error) {
			return analyzer.SearchCode(args)
		},
	)

	// Tool 6: Get file overview
	s.AddTool(
		mcp.NewTool("get_file_overview",
			mcp.WithDescription("Get a structured overview of a Go source file"),
			mcp.WithString("file_path",
				mcp.Description("Path to the Go file"),
				mcp.Required(),
			),
			mcp.WithString("level",
				mcp.Description("Detail level: outline (structure), signatures (with types), full (with docs)"),
				mcp.Enum("outline", "signatures", "full"),
			),
		),
		func(args map[string]interface{}) (*mcp.CallToolResult, error) {
			return analyzer.GetFileOverview(args)
		},
	)

	// Tool 7: Find dependencies
	s.AddTool(
		mcp.NewTool("find_dependencies",
			mcp.WithDescription("Find what packages/files depend on a given symbol or package"),
			mcp.WithString("target",
				mcp.Description("Package path or symbol name to find dependencies for"),
				mcp.Required(),
			),
			mcp.WithString("direction",
				mcp.Description("Dependency direction: forward (what it uses), reverse (what uses it), both"),
				mcp.Enum("forward", "reverse", "both"),
			),
		),
		func(args map[string]interface{}) (*mcp.CallToolResult, error) {
			return analyzer.FindDependencies(args)
		},
	)
}

func registerResources(s *server.MCPServer, analyzer *CodeAnalyzer) {
	// Resource template for package analysis
	s.AddResourceTemplate(
		mcp.ResourceTemplate{
			URITemplate: "code://package/{path}",
			Name:        "Go Package Analysis",
			Description: "Complete analysis of a Go package",
			MIMEType:    "application/json",
		},
		func(request mcp.ReadResourceRequest) ([]interface{}, error) {
			// Extract path from URI
			uri := request.Params.URI
			path := strings.TrimPrefix(uri, "code://package/")

			content, err := analyzer.GetPackageResource(path)
			if err != nil {
				return nil, err
			}

			return []interface{}{
				mcp.TextResourceContents{
					ResourceContents: mcp.ResourceContents{
						URI:      uri,
						MIMEType: "application/json",
					},
					Text: content,
				},
			}, nil
		},
	)

	// Resource template for file analysis
	s.AddResourceTemplate(
		mcp.ResourceTemplate{
			URITemplate: "code://file/{path}",
			Name:        "Go File Analysis",
			Description: "Complete analysis of a Go source file",
			MIMEType:    "application/json",
		},
		func(request mcp.ReadResourceRequest) ([]interface{}, error) {
			// Extract path from URI
			uri := request.Params.URI
			path := strings.TrimPrefix(uri, "code://file/")

			content, err := analyzer.GetFileResource(path)
			if err != nil {
				return nil, err
			}

			return []interface{}{
				mcp.TextResourceContents{
					ResourceContents: mcp.ResourceContents{
						URI:      uri,
						MIMEType: "application/json",
					},
					Text: content,
				},
			}, nil
		},
	)
}
