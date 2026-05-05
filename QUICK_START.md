# Quick Start Guide

Get up and running with Go AI Agent in 5 minutes!

## Prerequisites

1. **Go 1.21+** installed ([download](https://go.dev/dl/))
2. **Ollama** running locally ([install](https://ollama.ai))

## Step 1: Install Ollama Model

```bash
# Pull the Gemma 4 E2B model
ollama pull gemma4:e2b

# Verify it's available
ollama list
```

## Step 2: Build the Agent

```bash
# Navigate to project directory
cd /Users/ycchen/Desktop/agent

# Download dependencies
go mod download

# Build executable
go build -o agent ./cmd/agent

# Verify build
./agent -help
```

## Step 3: Run the Agent

```bash
# Start the agent
./agent

# You should see:
# ╔════════════════════════════════════════╗
# ║     Go AI Agent with Gemma 4 (E2B)     ║
# ╚════════════════════════════════════════╝
```

## Step 4: Try It Out

### Example 1: Basic Chat

```
[You]> Hello! What can you do?
[Gemma 4]> I'm an AI agent with several tools...
```

### Example 2: List Available Tools

```
[You]> /tools

Available Tools:
  [ENABLED] web_search
  [ENABLED] read_file
  [ENABLED] write_file
  [ENABLED] execute_system_command
  [ENABLED] calculator
  [ENABLED] extract_go_functions
  [ENABLED] extract_go_types
  [ENABLED] analyze_go_package
```

### Example 3: Use a Tool

```
[You]> Calculate 150 * 3.5 + 200 / 4
[Agent] Executing tool: calculator...
[Gemma 4]> The result is 575
```

### Example 4: Analyze Go Code

```
[You]> Analyze the structure of internal/agent/agent.go
[Agent] Executing tool: extract_go_types...
[Agent] Executing tool: extract_go_functions...
[Gemma 4]> The file contains:

Agent struct with fields:
  - modelName: string
  - ollamaClient: *ollama.Client
  - toolRegistry: *tools.Registry
  - permissions: *PermissionManager
  - conversation: *ConversationManager

Key functions:
  - NewAgent() - Creates a new agent instance
  - Chat() - Main conversation method
  - RegisterToolRegistry() - Registers tools
```

## Step 5: Explore Features

### Permission Control

```
[You]> /toggle execute_system_command
[System] Tool 'execute_system_command' is now DISABLED

[You]> Run ls command
[Gemma 4]> I cannot execute system commands as that tool is disabled.
```

### Clear History

```
[You]> /clear
[System] Conversation history cleared
```

### Exit

```
[You]> /exit
[System] Goodbye!
```

## Common Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/help` | Show help | `/help` |
| `/tools` | List all tools | `/tools` |
| `/toggle <tool>` | Enable/disable tool | `/toggle web_search` |
| `/clear` | Clear chat history | `/clear` |
| `/exit` | Quit agent | `/exit` |

## Example Use Cases

### 1. Code Analysis

```
Analyze the agent package and tell me what it does
```

### 2. Web Research

```
Search for Go best practices for error handling
```

### 3. File Operations

```
Read the README.md file and summarize it
```

### 4. Code Generation

```
Create a new Go file called utils.go with a string reversal function
```

### 5. Calculations

```
Calculate the fibonacci number at position 10
```

## Troubleshooting

### Problem: "connection refused"

**Solution:**
```bash
# Check Ollama is running
ollama list

# Restart Ollama if needed
# Then restart the agent
```

### Problem: "model not found"

**Solution:**
```bash
# Pull the model
ollama pull gemma4:e2b

# Verify it's available
ollama list
```

### Problem: "command not found: go"

**Solution:**
- Install Go from https://go.dev/dl/
- Add Go to your PATH

### Problem: Build errors

**Solution:**
```bash
# Clean and rebuild
go clean
go mod tidy
go build -o agent ./cmd/agent
```

## Next Steps

1. **Read full documentation:**
   - `README.md` - Complete overview
   - `USAGE_EXAMPLES.md` - More examples
   - `PROJECT_STATUS.md` - Current status

2. **Try advanced features:**
   - AST analysis on your own Go code
   - Multi-tool workflows
   - Permission management

3. **Customize:**
   - Use different Ollama models
   - Add your own tools
   - Configure via flags

## Tips for Best Results

1. **Be specific** in your requests
2. **Use AST tools** for code understanding instead of reading entire files
3. **Disable unused tools** for security
4. **Clear history** when switching contexts
5. **Check tool status** with `/tools` if something doesn't work

## Getting Help

- **Documentation:** Check README.md and other guides
- **Issues:** Report bugs on GitHub
- **Examples:** See USAGE_EXAMPLES.md

## Performance Tips

- Agent startup: ~10ms
- AST analysis: Fast (< 100ms)
- Web search: 2-5 seconds
- File operations: Near instant

## That's It!

You're ready to use the Go AI Agent. Have fun! 🚀

---

**Quick Reference:**

```bash
# Build
go build -o agent ./cmd/agent

# Run
./agent

# Help
./agent -help

# Custom model
./agent -model llama2:latest
```

For more details, see the full documentation in this directory.
