from typing import List, Callable, Dict
import ollama

class Agent:
    def __init__(self, model_name: str = "gemma4:4b"):
        """
        Initialize the AI Agent.
        
        Args:
            model_name: The name of the Ollama model to use. Default is gemma4:4b.
        """
        self.model_name = model_name
        self.tools: List[Callable] = []
        self.permissions: Dict[str, bool] = {}
        # Initialize conversation with system prompt
        self.messages = [
            {
                "role": "system", 
                "content": "You are a highly capable AI assistant powered by Gemma 4. "
                           "You have access to several tools. Always use them if they are needed to fulfill the user's request."
            }
        ]
        
    def register_tools(self, tools_list: List[Callable]):
        """Register tools that the agent can use."""
        self.tools = tools_list
        for tool in tools_list:
            self.permissions[tool.__name__] = True
            
    def toggle_permission(self, tool_name: str) -> str:
        """Toggle the execution permission for a specific tool."""
        if tool_name in self.permissions:
            self.permissions[tool_name] = not self.permissions[tool_name]
            state = "ENABLED" if self.permissions[tool_name] else "DISABLED"
            return f"Tool '{tool_name}' is now {state}."
        return f"Tool '{tool_name}' not found."
        
    def display_permissions(self) -> str:
        """Return a formatted string of current tool permissions."""
        res = "=== Tool Permissions ===\n"
        for name, is_enabled in self.permissions.items():
            state = "[O] ENABLED" if is_enabled else "[X] DISABLED"
            res += f"{state} : {name}\n"
        res += "========================\n"
        return res

    def chat(self, user_input: str):
        """
        Send user input to the agent and process any tool calls requested by the model.
        """
        self.messages.append({"role": "user", "content": user_input})
        
        while True:
            try:
                response = ollama.chat(
                    model=self.model_name,
                    messages=self.messages,
                    tools=self.tools
                )
            except Exception as e:
                return f"Error communicating with Ollama: {str(e)}\nPlease make sure Ollama is running and the model '{self.model_name}' is installed."
            
            message = response.get('message', {})
            self.messages.append(message)

            tool_calls = message.get('tool_calls')
            if not tool_calls:
                # No more tools to call, return the final response
                return message.get('content', '')
            
            # The model wants to call tools
            for tool_call in tool_calls:
                tool_name = tool_call['function']['name']
                tool_args = tool_call['function']['arguments']
                
                # Check user permissions
                if not self.permissions.get(tool_name, False):
                    print(f"\n[System] Blocked attempt to use tool: {tool_name}")
                    tool_result = f"Permission denied. The user has disabled the '{tool_name}' tool."
                else:
                    print(f"\n[Agent] Executing tool: {tool_name}({tool_args})...")
                    # Find and execute the tool
                    func = next((t for t in self.tools if t.__name__ == tool_name), None)
                    if func:
                        try:
                            # Pass dictionary arguments to the function
                            tool_result = str(func(**tool_args))
                        except Exception as e:
                            tool_result = f"Error executing tool: {e}"
                    else:
                        tool_result = f"Error: Tool {tool_name} not found."
                
                # Append tool execution result back to the conversation
                self.messages.append({
                    "role": "tool",
                    "content": tool_result,
                    "name": tool_name
                })
