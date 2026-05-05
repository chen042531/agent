# Integration Guide: Go Code Analyzer MCP with Claude

This guide shows how to integrate the Go Code Analyzer MCP server with Claude Desktop and Claude Code.

## For Claude Desktop

### Step 1: Build the MCP Server

```bash
cd /Users/ycchen/Desktop/agent/mcp-server
go mod download
go build -o go-code-mcp main.go analyzer.go
```

### Step 2: Locate Claude Desktop Config

The config file location depends on your OS:

**macOS:**
```
~/Library/Application Support/Claude/claude_desktop_config.json
```

**Windows:**
```
%APPDATA%\Claude\claude_desktop_config.json
```

**Linux:**
```
~/.config/Claude/claude_desktop_config.json
```

### Step 3: Add MCP Server Configuration

Edit the config file and add:

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

**Important:** Update paths to match your system!

### Step 4: Restart Claude Desktop

Completely quit and restart Claude Desktop for changes to take effect.

### Step 5: Verify Integration

In Claude Desktop, start a conversation:

```
You: "What MCP servers are available?"
Claude: "I have access to the go-code-analyzer MCP server..."

You: "List all packages in this Go project"
Claude: [Uses list_packages tool automatically]
```

## For Claude Code (CLI)

### Step 1: Build the Server

Same as Claude Desktop:

```bash
cd /Users/ycchen/Desktop/agent/mcp-server
go build -o go-code-mcp main.go analyzer.go
```

### Step 2: Configure MCP Settings

Create or edit `~/.claude/mcp.json`:

```json
{
  "servers": {
    "go-code-analyzer": {
      "command": "/Users/ycchen/Desktop/agent/mcp-server/go-code-mcp",
      "env": {
        "WORKSPACE_ROOT": "${workspaceFolder}"
      }
    }
  }
}
```

The `${workspaceFolder}` variable will be replaced with the current working directory.

### Step 3: Use in Claude Code

```bash
# Navigate to your Go project
cd /Users/ycchen/Desktop/agent

# Start Claude Code
claude

# Claude will automatically have access to MCP tools
You: "Analyze the internal/agent package"
Claude: [Uses analyze_package tool]
```

## Dynamic Context Example

Here's how the MCP server provides context dynamically:

### Scenario: Understanding a New Codebase

```
You: "I'm new to this codebase, give me an overview"

Claude's workflow (automatic):

1. Tool: list_packages(level: "summary")
   Result: {
     packages: [
       {path: "internal/agent", type_count: 3, function_count: 12},
       {path: "internal/tools", type_count: 8, function_count: 25},
       ...
     ]
   }

2. Tool: analyze_package(
     package_path: "internal/agent",
     include: ["types", "functions"],
     exported_only: true
   )
   Result: {
     types: [
       {name: "Agent", kind: "struct", line: 15},
       {name: "PermissionManager", kind: "struct", line: 9}
     ],
     functions: [
       {name: "NewAgent", signature: "func(...) (*Agent, error)"},
       {name: "Chat", signature: "func(...) (string, error)"}
     ]
   }

3. Tool: get_type_details(
     file_path: "internal/agent/agent.go",
     type_name: "Agent",
     include_methods: true
   )
   Result: {
     name: "Agent",
     kind: "struct",
     fields: [...],
     methods: [...]
   }

Claude then synthesizes: "This is an AI agent system with the following architecture..."
```

### Scenario: Implementing a Feature

```
You: "How do I add a new tool to the agent?"

Claude's workflow:

1. Tool: search_code(query: "Tool", kind: "interface")
   Result: [
     {kind: "interface", name: "Tool", file: "internal/tools/types.go", line: 9}
   ]

2. Tool: get_type_details(
     file_path: "internal/tools/types.go",
     type_name: "Tool",
     include_methods: true
   )
   Result: {
     kind: "interface",
     interface_methods: [
       {name: "Name", signature: "func() string"},
       {name: "Description", signature: "func() string"},
       {name: "Parameters", signature: "func() ParameterSchema"},
       {name: "Execute", signature: "func(context.Context, map[string]interface{}) (string, error)"}
     ]
   }

3. Tool: search_code(query: "Tool", kind: "type")
   Result: [
     {name: "WebSearchTool", file: "internal/tools/web_search.go"},
     {name: "ReadFileTool", file: "internal/tools/file_operations.go"},
     ...
   ]

4. Tool: get_function_details(
     file_path: "internal/tools/web_search.go",
     function_name: "WebSearchTool.Execute",
     include_body: true
   )
   Result: {full implementation}

Claude then explains: "To add a new tool, implement the Tool interface..."
```

## Progressive Depth Control

The MCP server allows Claude to control context depth:

### Level 1: Minimal (Names Only)

```json
{
  "tool": "list_packages",
  "args": {"level": "minimal"}
}
```

Returns: `["agent", "cli", "ollama", "tools"]`

**Use case:** Quick discovery, low token usage

### Level 2: Summary (Counts)

```json
{
  "tool": "list_packages",
  "args": {"level": "summary"}
}
```

Returns:
```json
[
  {"path": "agent", "type_count": 3, "function_count": 12},
  {"path": "cli", "type_count": 2, "function_count": 8}
]
```

**Use case:** Understanding scope

### Level 3: Detailed (Structure)

```json
{
  "tool": "analyze_package",
  "args": {
    "package_path": "internal/agent",
    "include": ["types", "functions"],
    "exported_only": true
  }
}
```

Returns: Type and function signatures

**Use case:** Understanding API

### Level 4: Full (With Bodies)

```json
{
  "tool": "get_function_details",
  "args": {
    "file_path": "internal/agent/agent.go",
    "function_name": "Chat",
    "include_body": true
  }
}
```

Returns: Complete implementation

**Use case:** Deep understanding, debugging

## Environment Configuration

### For Different Projects

You can configure multiple MCP servers for different projects:

```json
{
  "mcpServers": {
    "agent-project": {
      "command": "/path/to/go-code-mcp",
      "env": {
        "WORKSPACE_ROOT": "/Users/ycchen/Desktop/agent"
      }
    },
    "other-project": {
      "command": "/path/to/go-code-mcp",
      "env": {
        "WORKSPACE_ROOT": "/path/to/other/project"
      }
    }
  }
}
```

### With Logging

Enable detailed logging:

```json
{
  "mcpServers": {
    "go-code-analyzer": {
      "command": "/path/to/go-code-mcp",
      "env": {
        "WORKSPACE_ROOT": "/Users/ycchen/Desktop/agent",
        "MCP_LOG_LEVEL": "debug"
      }
    }
  }
}
```

View logs:
```bash
tail -f ~/Library/Logs/Claude/mcp-server-go-code-analyzer.log
```

## Verification

### Test 1: Check Server Availability

In Claude:
```
You: "What MCP tools do you have access to?"
```

Expected: Claude lists the 7 Go analysis tools

### Test 2: Basic Function

In Claude:
```
You: "List all packages in this project"
```

Expected: Claude uses `list_packages` and shows results

### Test 3: Deep Analysis

In Claude:
```
You: "Explain the Agent struct in detail"
```

Expected: Claude uses multiple tools to build complete picture

## Troubleshooting

### Issue: Tools Not Appearing

**Check 1:** Config file syntax
```bash
# Validate JSON
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json | jq .
```

**Check 2:** Server executable
```bash
# Test server directly
/path/to/go-code-mcp
# Should start and wait for stdio input
```

**Check 3:** Logs
```bash
# macOS
cat ~/Library/Logs/Claude/mcp-server-go-code-analyzer.log

# Windows
type %APPDATA%\Claude\Logs\mcp-server-go-code-analyzer.log
```

### Issue: Server Crashes

**Check:** WORKSPACE_ROOT exists and has Go code
```bash
ls -la $WORKSPACE_ROOT
```

**Check:** Go files are valid
```bash
cd $WORKSPACE_ROOT
go build ./...
```

### Issue: Slow Performance

**Solution 1:** Use exported_only=true
```json
{"exported_only": true}
```

**Solution 2:** Narrow scope
```json
{"package_path": "internal/agent"}  // Not "."
```

**Solution 3:** Use minimal level first
```json
{"level": "minimal"}
```

## Best Practices

### 1. Start Broad, Then Narrow

❌ Bad:
```
"Show me all code in the project"
```

✅ Good:
```
"List all packages" → "Analyze the agent package" → "Show me the Chat method"
```

### 2. Use Appropriate Detail Levels

❌ Bad: Always use "full"
✅ Good: Use "minimal" → "summary" → "detailed" → "full" as needed

### 3. Filter Early

❌ Bad:
```json
{"exported_only": false, "include": ["everything"]}
```

✅ Good:
```json
{"exported_only": true, "include": ["types", "functions"]}
```

### 4. Leverage Search

❌ Bad: Analyze entire codebase to find one function
✅ Good: Use search_code first, then get details

## Advanced Usage

### Custom Workflows

You can create custom workflows by chaining tools:

**Workflow: "Find similar code"**
1. `search_code` - Find candidates
2. `get_function_details` - Get implementations
3. Compare and analyze

**Workflow: "Understand dependencies"**
1. `analyze_package` - Get imports
2. `search_code` - Find where used
3. `find_dependencies` - Build graph

**Workflow: "Refactoring analysis"**
1. `get_type_details` - Understand type
2. `search_code` - Find all usage
3. `get_function_details` - Check each usage

## Integration with Development Workflow

### During Code Review

```
Reviewer: "Claude, analyze the changes in this PR"
Claude: [Uses MCP to understand current code]
Claude: "I see these changes affect the Permission system..."
```

### During Development

```
Developer: "Claude, help me implement feature X"
Claude: [Uses MCP to understand existing code]
Claude: "Based on the existing architecture, here's how..."
```

### During Onboarding

```
New Developer: "Claude, explain this codebase"
Claude: [Progressive exploration via MCP]
Claude: "This is a Go AI agent with the following components..."
```

## Conclusion

The Go Code Analyzer MCP server enables Claude to:

✅ Explore codebases progressively
✅ Adjust context depth dynamically
✅ Minimize token usage while maximizing relevance
✅ Provide intelligent, context-aware assistance

Start using it today to supercharge your development workflow with Claude!
