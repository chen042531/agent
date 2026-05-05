package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ycchen/ai-agent-go/internal/agent"
)

type REPL struct {
	agent          *agent.Agent
	commandHandler *CommandHandler
	running        bool
}

func NewREPL(agent *agent.Agent) *REPL {
	return &REPL{
		agent:          agent,
		commandHandler: NewCommandHandler(agent),
		running:        false,
	}
}

func (r *REPL) Run(ctx context.Context) error {
	r.running = true
	PrintWelcome()

	scanner := bufio.NewScanner(os.Stdin)

	for r.running {
		PrintUserPrompt()

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Check if it's a command
		isCommand, err := r.commandHandler.HandleCommand(input)
		if isCommand {
			if err != nil {
				if err.Error() == "exit" {
					r.running = false
					break
				}
				PrintError(err)
			}
			continue
		}

		// Regular chat
		response, err := r.agent.Chat(ctx, input)
		if err != nil {
			PrintError(err)
			continue
		}

		PrintAssistantMessage(response)
		PrintSeparator()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}
