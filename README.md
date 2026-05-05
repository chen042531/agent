# Go AI Agent with Gemma 4 (E2B)

A powerful AI agent built in Go that integrates with Ollama's Gemma 4 model. Features include web search, file operations, system commands, mathematical calculations, and advanced Go code AST analysis.

## Features

### Core Capabilities
- **Conversational AI**: Natural language interaction powered by Gemma 4 (E2B) via Ollama
- **Tool Execution**: Dynamic tool calling with permission management
- **Conversation History**: Maintains context across multiple turns

### Built-in Tools

1. **Web Search** (`web_search`)
   - Search the web using DuckDuckGo
   - Returns top 5 results with titles, snippets, and URLs

2. **File Operations**
   - `read_file`: Read file contents
   - `write_file`: Write/create files with automatic directory creation

3. **System Commands** (`execute_system_command`)
   - Execute shell commands with 30-second timeout
   - Security: sandboxed execution with output capture

4. **Calculator** (`calculator`)
   - Safe mathematical expression evaluation
   - Supports standard operators: +, -, *, /, ()

5. **Go AST Analysis** (New!)
   - `extract_go_functions`: Extract all function signatures from Go files
   - `extract_go_types`: Extract type definitions (structs, interfaces, aliases)
   - `analyze_go_package`: Analyze entire Go packages for exported symbols

## Prerequisites

- **Go**: 1.21 or higher
- **Ollama**: Running locally with Gemma 4 E2B model
  ```bash
  # Install Ollama: https://ollama.ai
  # Pull the model
  ollama pull gemma4:e2b
  ```

## Installation

```bash
cd /Users/ycchen/Desktop/agent

# Install dependencies
go mod download

# Build the agent
go build -o agent ./cmd/agent

# Or run directly
go run ./cmd/agent/main.go
```

## Usage

### Basic Usage

```bash
# Start the agent
./agent

# Or with custom settings
./agent -model gemma4:e2b -ollama-url http://localhost:11434
```

### CLI Commands

- `/help` - Show help message
- `/tools` - List all tools and their permission status
- `/toggle <tool_name>` - Enable/disable a specific tool
- `/clear` - Clear conversation history
- `/exit` - Exit the program

## Project Structure

```
agent/
├── cmd/agent/main.go
├── internal/
│   ├── agent/         # Core agent logic
│   ├── ollama/        # Ollama API client
│   ├── tools/         # Tool implementations
│   │   └── ast/       # Go AST analysis
│   └── cli/           # CLI interface
└── go.mod
```

## License

MIT License
