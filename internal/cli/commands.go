package cli

import (
	"fmt"
	"strings"

	"github.com/ycchen/ai-agent-go/internal/agent"
)

type CommandHandler struct {
	agent *agent.Agent
}

func NewCommandHandler(agent *agent.Agent) *CommandHandler {
	return &CommandHandler{
		agent: agent,
	}
}

func (h *CommandHandler) HandleCommand(input string) (bool, error) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "/") {
		return false, nil // Not a command
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false, nil
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "/help":
		PrintHelp()
		return true, nil

	case "/tools":
		h.listTools()
		return true, nil

	case "/toggle":
		if len(args) < 1 {
			return true, fmt.Errorf("usage: /toggle <tool_name>")
		}
		h.toggleTool(args[0])
		return true, nil

	case "/clear":
		h.agent.ClearConversation()
		PrintSystemMessage("Conversation history cleared")
		return true, nil

	case "/exit":
		PrintSystemMessage("Goodbye!")
		return true, fmt.Errorf("exit")

	default:
		return true, fmt.Errorf("unknown command: %s (type /help for available commands)", command)
	}
}

func (h *CommandHandler) listTools() {
	registry := h.agent.GetToolRegistry()
	permissions := h.agent.GetPermissionManager()

	tools := registry.All()
	if len(tools) == 0 {
		PrintSystemMessage("No tools registered")
		return
	}

	fmt.Println(ColorYellow + "\nAvailable Tools:" + ColorReset)
	for _, tool := range tools {
		status := "ENABLED"
		statusColor := ColorGreen
		if !permissions.IsAllowed(tool.Name()) {
			status = "DISABLED"
			statusColor = ColorRed
		}

		fmt.Printf("  %s[%s]%s %s\n", statusColor, status, ColorReset, tool.Name())
		fmt.Printf("    %s\n", tool.Description())
	}
	fmt.Println()
}

func (h *CommandHandler) toggleTool(toolName string) {
	registry := h.agent.GetToolRegistry()
	permissions := h.agent.GetPermissionManager()

	// Check if tool exists
	_, exists := registry.Get(toolName)
	if !exists {
		PrintError(fmt.Errorf("tool '%s' not found", toolName))
		return
	}

	permissions.Toggle(toolName)
	newStatus := "ENABLED"
	statusColor := ColorGreen
	if !permissions.IsAllowed(toolName) {
		newStatus = "DISABLED"
		statusColor = ColorRed
	}

	fmt.Printf("%s[System] Tool '%s' is now %s%s%s%s\n",
		ColorYellow, toolName, statusColor, newStatus, ColorReset, "\n")
}
