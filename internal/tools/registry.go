package tools

import (
	"sync"

	"github.com/ycchen/ai-agent-go/internal/ollama"
)

type Registry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(tool Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[tool.Name()] = tool
}

func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, ok := r.tools[name]
	return tool, ok
}

func (r *Registry) All() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

func (r *Registry) ToOllamaTools() []ollama.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ollamaTools := make([]ollama.Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		ollamaTools = append(ollamaTools, ollama.Tool{
			Type: "function",
			Function: ollama.FunctionSpec{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters:  tool.Parameters(),
			},
		})
	}
	return ollamaTools
}
