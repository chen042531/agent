# Usage Examples

## Quick Start

```bash
# Build and run
cd /Users/ycchen/Desktop/agent
go build -o agent ./cmd/agent
./agent
```

## Example 1: Analyzing Go Code Structure

**User Input:**
```
Analyze test_sample.go and tell me what types and functions it contains
```

**Expected Tool Calls:**
1. `extract_go_types` on test_sample.go
2. `extract_go_functions` on test_sample.go

**Expected Output:**
```
The file contains:

Type: User (struct) at line 6
  Fields:
    - ID int (json:"id")
    - Name string (json:"name")
    - Email string (json:"email")
    - IsActive bool (json:"is_active")

Functions:
  - GetUser(id int) (*User, error) at line 13
  - (User) UpdateEmail(email string) error at line 20
  - main() at line 26
```

## Example 2: Package-Level Analysis

**User Input:**
```
Analyze the internal/agent package structure
```

**Expected Tool Calls:**
1. `analyze_go_package` on internal/agent

**Expected Output:**
```
Package: agent

Exported Types:
  - Agent (struct)
  - PermissionManager (struct)
  - ConversationManager (struct)

Exported Functions:
  - NewAgent(modelName string, ollamaURL string) (*Agent, error)
  - NewPermissionManager() *PermissionManager
  - NewConversationManager() *ConversationManager
  - (Agent) Chat(ctx context.Context, userInput string) (string, error)
  - (PermissionManager) IsAllowed(toolName string) bool
  - (PermissionManager) Toggle(toolName string)
```

## Example 3: Code Generation with Context

**User Input:**
```
Read internal/agent/agent.go and add a new method to save conversation history to a file
```

**Expected Tool Calls:**
1. `read_file` on internal/agent/agent.go
2. `extract_go_types` to understand Agent struct
3. `write_file` to add the new method

**Agent Actions:**
- Understands existing Agent structure from AST
- Generates appropriate Go code following conventions
- Adds method with proper error handling

## Example 4: Web Search + Code Analysis

**User Input:**
```
Search for best practices on Go error handling and check if our agent.go follows them
```

**Expected Tool Calls:**
1. `web_search` for "Go error handling best practices"
2. `read_file` on internal/agent/agent.go
3. `extract_go_functions` to see function signatures

**Expected Output:**
```
Based on Go best practices:
✓ Functions return errors as last return value
✓ Error wrapping with fmt.Errorf and %w
✓ Context passed as first parameter

Areas for improvement:
- Consider adding custom error types for specific cases
- Add error logging in some paths
```

## Example 5: Permission Control

**Session:**
```
[You]> /tools
[System] Available Tools:
  [ENABLED] web_search
  [ENABLED] execute_system_command
  [ENABLED] read_file
  ...

[You]> /toggle execute_system_command
[System] Tool 'execute_system_command' is now DISABLED

[You]> Run ls command
[Gemma 4]> I apologize, but I cannot execute system commands as that tool is currently disabled.

[You]> /toggle execute_system_command
[System] Tool 'execute_system_command' is now ENABLED

[You]> Run ls command
[Agent] Executing tool: execute_system_command...
[Gemma 4]> Here are the files in the current directory:
agent
cmd
go.mod
go.sum
internal
README.md
```

## Example 6: Calculator Tool

**User Input:**
```
Calculate (150 * 3.5) + (200 / 4)
```

**Expected Tool Calls:**
1. `calculator` with expression "(150 * 3.5) + (200 / 4)"

**Expected Output:**
```
The result is 575
```

## Example 7: Multi-Tool Workflow

**User Input:**
```
Create a new Go file called utils.go with a function that validates email addresses
```

**Expected Tool Calls:**
1. `web_search` (optional) - find email regex pattern
2. `write_file` - create utils.go with validation function

**Expected Output:**
```go
package main

import "regexp"

// ValidateEmail checks if an email address is valid
func ValidateEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(pattern)
    return re.MatchString(email)
}
```

## Example 8: Refactoring Assistant

**User Input:**
```
Read internal/tools/registry.go and suggest how to add a method to list tools by category
```

**Expected Tool Calls:**
1. `read_file` on internal/tools/registry.go
2. `extract_go_types` to see Registry struct
3. `extract_go_functions` to see existing methods

**Agent Actions:**
- Analyzes current structure
- Suggests adding a `Category` field to Tool interface
- Proposes `ListByCategory()` method implementation

## Tips for Best Results

1. **Be Specific**: Instead of "analyze this", say "extract functions from internal/agent/agent.go"

2. **Combine Tools**: The agent can chain multiple tools automatically
   - "Search for X and create a file with the results"
   - "Read file.go and explain its structure"

3. **Use AST Tools for Code Understanding**:
   - Before modifying code, use AST tools to understand structure
   - This gives the LLM precise context instead of full source

4. **Permission Management**:
   - Disable risky tools when not needed
   - Use `/toggle` to control access

5. **Clear Conversations**:
   - Use `/clear` when switching contexts
   - Prevents confusion from old conversation history

## Performance Notes

- **AST Analysis**: Very fast (< 100ms for most files)
- **Web Search**: 2-5 seconds depending on network
- **File Operations**: Instant for files < 1MB
- **System Commands**: Up to 30 seconds timeout
- **Calculator**: Instant evaluation

## Common Patterns

### Pattern 1: Code Exploration
```
1. "What's in package X?" → analyze_go_package
2. "Show me the Agent struct" → extract_go_types
3. "What methods does Agent have?" → extract_go_functions
```

### Pattern 2: Code Generation
```
1. Read existing code for context
2. Extract types/functions to understand structure
3. Generate new code following conventions
4. Write to file
```

### Pattern 3: Research + Implementation
```
1. Web search for information/examples
2. Extract relevant patterns
3. Apply to current codebase
4. Verify with AST analysis
```

## Troubleshooting

**Agent doesn't call AST tools:**
- Be explicit: "use extract_go_functions on agent.go"
- Or: "analyze the structure of agent.go using AST"

**Tool disabled unexpectedly:**
- Check with `/tools`
- Re-enable with `/toggle <tool_name>`

**Ollama connection issues:**
- Ensure Ollama is running: `ollama list`
- Check model is available: `ollama pull gemma4:e2b`
- Verify URL: default is http://localhost:11434
