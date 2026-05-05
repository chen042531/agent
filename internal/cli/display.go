package cli

import (
	"fmt"
	"strings"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

func PrintWelcome() {
	fmt.Println(ColorBold + ColorCyan + "╔════════════════════════════════════════╗" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "║     Go AI Agent with Gemma 4 (E2B)     ║" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "╚════════════════════════════════════════╝" + ColorReset)
	fmt.Println()
	fmt.Println(ColorYellow + "Commands:" + ColorReset)
	fmt.Println("  /help      - Show this help message")
	fmt.Println("  /tools     - List all available tools")
	fmt.Println("  /toggle    - Toggle tool permission (e.g., /toggle web_search)")
	fmt.Println("  /clear     - Clear conversation history")
	fmt.Println("  /exit      - Exit the program")
	fmt.Println()
	fmt.Println(ColorGreen + "Type your message and press Enter to chat!" + ColorReset)
	fmt.Println(strings.Repeat("-", 60))
}

func PrintHelp() {
	fmt.Println(ColorYellow + "\nAvailable Commands:" + ColorReset)
	fmt.Println("  /help      - Show this help message")
	fmt.Println("  /tools     - List all available tools and their status")
	fmt.Println("  /toggle    - Toggle tool permission")
	fmt.Println("               Example: /toggle web_search")
	fmt.Println("  /clear     - Clear conversation history")
	fmt.Println("  /exit      - Exit the agent")
	fmt.Println()
}

func PrintUserPrompt() {
	fmt.Print(ColorBold + ColorBlue + "[You]> " + ColorReset)
}

func PrintAssistantMessage(msg string) {
	fmt.Println(ColorBold + ColorGreen + "[Gemma 4]> " + ColorReset + msg)
	fmt.Println()
}

func PrintSystemMessage(msg string) {
	fmt.Println(ColorYellow + "[System] " + ColorReset + msg)
}

func PrintError(err error) {
	fmt.Println(ColorRed + "[Error] " + ColorReset + err.Error())
}

func PrintToolExecution(toolName string) {
	fmt.Println(ColorPurple + "[Agent] " + ColorReset + "Executing tool: " + ColorCyan + toolName + ColorReset + "...")
}

func PrintSeparator() {
	fmt.Println(strings.Repeat("-", 60))
}
