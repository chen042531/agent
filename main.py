import sys
from agent import Agent
from tools import AVAILABLE_TOOLS

def print_help():
    print("""
=== Gemma 4 Agent Commands ===
/help             - Show this help message
/tools            - List all tools and their current permission status
/toggle <tool>    - Enable or disable a specific tool (e.g. /toggle web_search)
/clear            - Clear the conversation history
/exit             - Exit the application
==============================
""")

def main():
    print("Welcome to Gemma 4 AI Agent!")
    print("Initializing agent...")
    
    # You can change to gemma4:26b or gemma4:31b if you have enough RAM
    agent = Agent(model_name="gemma4:4b")
    agent.register_tools(AVAILABLE_TOOLS)
    
    print("Agent initialized successfully.")
    print_help()
    
    while True:
        try:
            user_input = input("\n[You]> ").strip()
        except (KeyboardInterrupt, EOFError):
            print("\nExiting...")
            break
            
        if not user_input:
            continue
            
        if user_input.lower() == '/exit':
            print("Goodbye!")
            break
        elif user_input.lower() == '/help':
            print_help()
        elif user_input.lower() == '/tools':
            print(agent.display_permissions())
        elif user_input.lower().startswith('/toggle '):
            tool_name = user_input.split(' ', 1)[1].strip()
            result = agent.toggle_permission(tool_name)
            print(f"[System] {result}")
        elif user_input.lower() == '/clear':
            # Reset conversation (keeping the system prompt)
            agent.messages = [agent.messages[0]]
            print("[System] Conversation history cleared.")
        else:
            # Send to agent
            print("[Agent is thinking...]")
            response = agent.chat(user_input)
            print(f"\n[Gemma 4]> {response}")

if __name__ == "__main__":
    main()
