

# Go Code Analyzer MCP Server

A Model Context Protocol (MCP) server that provides intelligent, context-aware Go code analysis for Claude and other LLM applications.

## 🌟 Key Features

### Dynamic Context Adjustment
- **Progressive Exploration**: Start with high-level summaries, drill down as needed
- **Configurable Depth**: Choose from minimal, summary, detailed, or full analysis
- **Smart Filtering**: Export-only, specific kinds, or full codebase analysis

### Intelligent Code Analysis Tools

1. **`list_packages`** - Discover packages with configurable detail levels
2. **`analyze_package`** - Deep package analysis with selective includes
3. **`get_function_details`** - Detailed function/method information
4. **`get_type_details`** - Comprehensive type analysis with methods
5. **`search_code`** - Fast symbol search across the codebase
6. **`get_file_overview`** - Structured file summaries
7. **`find_dependencies`** - Dependency graph analysis (coming soon)

### Context-Aware Design

The server is designed to work seamlessly with Claude's iterative reasoning:

```
Claude: "What packages exist here?"
→ Tool: list_packages (level: minimal)
→ Response: ["agent", "cli", "ollama", "tools"]

Claude: "Tell me more about the tools package"
→ Tool: analyze_package (include: ["types", "functions"], exported_only: true)
→ Response: {types: [...], functions: [...]}

Claude: "Show me the WebSearchTool implementation details"
→ Tool: get_type_details (include_methods: true)
→ Response: {fields: [...], methods: [...], documentation: "..."}
```

## 🚀 Quick Start

### 1. Build the Server

```bash
cd mcp-server
go mod download
go build -o go-code-mcp
```

### 2. Configure Claude Desktop

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "go-code-analyzer": {
      "command": "/Users/ycchen/Desktop/agent/mcp-server/go-code-mcp",
      "env": {
        "WORKSPACE_ROOT": "/Users/ycchen/Desktop/agent"
      }
    }
  }
}
```

### 3. Restart Claude Desktop

The MCP server will be available as a tool in Claude Desktop.

## 📖 Tool Reference

### list_packages

Lists all Go packages in the workspace.

**Parameters:**
- `path` (string, optional): Relative path from workspace root
- `level` (string, optional): Detail level
  - `minimal`: Package paths only
  - `summary`: Paths + type/function counts (default)
  - `detailed`: Paths + counts + exported symbols list

**Example:**
```
List all packages with summary info
```

Claude will call:
```json
{
  "path": ".",
  "level": "summary"
}
```

### analyze_package

Deep analysis of a specific package.

**Parameters:**
- `package_path` (string, required): Path to package directory
- `include` (array, optional): What to include
  - Options: `types`, `functions`, `imports`, `methods`, `fields`
  - Default: `["types", "functions"]`
- `exported_only` (boolean, optional): Only exported symbols (default: true)

**Example:**
```
Analyze the internal/agent package, include all details
```

Claude will call:
```json
{
  "package_path": "internal/agent",
  "include": ["types", "functions", "imports", "methods", "fields"],
  "exported_only": false
}
```

### get_function_details

Detailed information about a specific function.

**Parameters:**
- `file_path` (string, required): Path to Go file
- `function_name` (string, required): Function name (use `Type.Method` for methods)
- `include_body` (boolean, optional): Include source code (default: false)

**Example:**
```
Show me the implementation of Agent.Chat method
```

Claude will call:
```json
{
  "file_path": "internal/agent/agent.go",
  "function_name": "Agent.Chat",
  "include_body": true
}
```

### get_type_details

Comprehensive type information.

**Parameters:**
- `file_path` (string, required): Path to Go file
- `type_name` (string, required): Type name
- `include_methods` (boolean, optional): Include type methods (default: true)

**Example:**
```
What is the Agent struct and what methods does it have?
```

Claude will call:
```json
{
  "file_path": "internal/agent/agent.go",
  "type_name": "Agent",
  "include_methods": true
}
```

### search_code

Search for symbols across the codebase.

**Parameters:**
- `query` (string, required): Search query (name pattern)
- `kind` (string, optional): Symbol type
  - Options: `all`, `function`, `type`, `method`, `interface`
  - Default: `all`
- `max_results` (integer, optional): Max results (default: 20)

**Example:**
```
Find all functions related to "Parse"
```

Claude will call:
```json
{
  "query": "Parse",
  "kind": "function",
  "max_results": 10
}
```

### get_file_overview

Structured overview of a Go file.

**Parameters:**
- `file_path` (string, required): Path to Go file
- `level` (string, optional): Detail level
  - `outline`: Structure only (names, line numbers)
  - `signatures`: With type signatures (default)
  - `full`: With documentation

**Example:**
```
Give me an overview of the main.go file
```

Claude will call:
```json
{
  "file_path": "cmd/agent/main.go",
  "level": "full"
}
```

## 🎯 Usage Patterns

### Pattern 1: Exploratory Analysis

```
Human: "Help me understand this codebase"

Claude uses:
1. list_packages (level: summary) → Get overview
2. analyze_package (for key packages) → Understand structure
3. get_type_details (for main types) → Deep dive
```

### Pattern 2: Feature Implementation

```
Human: "How would I add a new tool?"

Claude uses:
1. search_code (query: "Tool") → Find Tool interface
2. get_type_details (type_name: "Tool") → Understand interface
3. get_function_details (example tool) → See implementation pattern
```

### Pattern 3: Debugging

```
Human: "Why isn't the permission system working?"

Claude uses:
1. search_code (query: "permission") → Find related code
2. get_function_details (include_body: true) → See implementation
3. analyze_package (include: ["functions", "types"]) → Understand flow
```

### Pattern 4: Refactoring

```
Human: "Should we split the agent package?"

Claude uses:
1. analyze_package (include: all) → Full package analysis
2. find_dependencies (direction: "both") → See coupling
3. get_file_overview (for each file) → Understand cohesion
```

## 🧠 How It Works with Claude

### Context Optimization

The MCP server is designed to minimize token usage while maximizing relevance:

1. **Lazy Loading**: Claude starts with summaries and requests details only when needed
2. **Selective Fields**: Include only relevant information per query
3. **Smart Defaults**: Reasonable defaults for common use cases
4. **Progressive Depth**: Move from overview → details → source code

### Example Flow

```
User: "Add logging to the Chat method"

Step 1: Find the method
Claude → search_code(query: "Chat", kind: "method")
Server → [{name: "Chat", file: "internal/agent/agent.go", line: 34}]

Step 2: Understand current implementation
Claude → get_function_details(
  file_path: "internal/agent/agent.go",
  function_name: "Agent.Chat",
  include_body: true
)
Server → {signature: "...", body: "...", documentation: "..."}

Step 3: Check if logging already exists
Claude analyzes body, sees no logging

Step 4: Check what logging is used elsewhere
Claude → search_code(query: "log", kind: "all")
Server → [import paths, usage examples]

Step 5: Implement
Claude suggests code changes with proper logging
```

## 🔧 Advanced Configuration

### Environment Variables

- `WORKSPACE_ROOT`: Root directory to analyze (default: current directory)
- `MCP_LOG_LEVEL`: Logging level (default: info)

### Custom Filters

The analyzer supports custom filtering:

```json
{
  "include": ["types", "functions"],
  "exported_only": true  // Only public API
}
```

vs.

```json
{
  "include": ["types", "functions", "methods", "fields"],
  "exported_only": false  // Everything
}
```

## 📊 Performance

- **Fast Parsing**: Uses Go's built-in `go/ast` package
- **No External Dependencies**: Pure standard library for analysis
- **Incremental Analysis**: Only parses files/packages as needed
- **Caching**: (Coming soon) Cache parsed ASTs for repeated queries

## 🛠️ Development

### Running Locally

```bash
# Set workspace
export WORKSPACE_ROOT=/path/to/your/go/project

# Run server
go run main.go analyzer.go
```

### Testing with MCP Inspector

```bash
# Install MCP inspector
npm install -g @modelcontextprotocol/inspector

# Run inspector
mcp-inspector ./go-code-mcp
```

### Debugging

The server logs to stderr. To see logs:

```bash
./go-code-mcp 2> mcp-server.log
```

## 🚧 Roadmap

### v1.1 (Current)
- [x] Core MCP server
- [x] 7 analysis tools
- [x] Progressive depth control
- [ ] Resource templates
- [ ] Comprehensive tests

### v1.2 (Planned)
- [ ] Dependency graph analysis
- [ ] Call graph visualization
- [ ] Cross-package analysis
- [ ] AST caching
- [ ] Performance metrics

### v2.0 (Future)
- [ ] Multi-language support (Python, TypeScript)
- [ ] Semantic code search
- [ ] Code generation helpers
- [ ] Refactoring suggestions
- [ ] Test generation

## 📝 Examples

### Example 1: Understanding a New Codebase

```
User: "Explain how this Go agent works"

Claude's approach:
1. list_packages → See structure
2. analyze_package("internal/agent") → Understand core
3. get_type_details("Agent") → See main type
4. get_function_details("Agent.Chat") → Understand flow
5. Explain architecture to user
```

### Example 2: Adding a Feature

```
User: "Add a tool to read JSON files"

Claude's approach:
1. search_code("Tool") → Find tool interface
2. get_type_details("Tool") → Understand interface
3. search_code("ReadFile") → Find similar tool
4. get_function_details("ReadFileTool.Execute", include_body: true)
5. Generate new ReadJSONTool implementation
```

### Example 3: Debugging

```
User: "Why do my tool calls fail?"

Claude's approach:
1. search_code("Execute") → Find execution paths
2. analyze_package("internal/tools") → See tool registry
3. get_function_details("Agent.Chat", include_body: true)
4. Identify permission check logic
5. Explain issue and solution
```

## 🤝 Integration with Claude Desktop

When properly configured, Claude Desktop will:

1. **Automatically discover** available MCP tools
2. **Intelligently choose** the right tool for each query
3. **Chain tool calls** to build understanding progressively
4. **Optimize context** by requesting only needed details

The result is a natural conversation where Claude explores your codebase like an experienced developer.

## 📚 Learn More

- [MCP Specification](https://modelcontextprotocol.io/)
- [Claude Desktop Configuration](https://docs.anthropic.com/claude/docs/claude-desktop)
- [Go AST Package](https://pkg.go.dev/go/ast)

## 🐛 Troubleshooting

**Server doesn't start:**
- Check WORKSPACE_ROOT is valid
- Ensure Go code is valid in workspace
- Check Claude Desktop logs

**No tools appear in Claude:**
- Verify config path is correct
- Restart Claude Desktop completely
- Check server is executable

**Slow responses:**
- Large codebases may take time to parse
- Consider narrowing analysis scope
- Use exported_only=true when possible

**Permission errors:**
- Ensure server has read access to WORKSPACE_ROOT
- Check file permissions

## 🎉 Conclusion

This MCP server transforms how Claude interacts with Go codebases, enabling:

- **Efficient exploration** through progressive detail
- **Intelligent analysis** with context-aware tools
- **Natural workflow** that matches developer thinking
- **Minimal overhead** with smart filtering

Start exploring your Go codebase with Claude today!
