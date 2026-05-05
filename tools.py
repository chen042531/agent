import os
import subprocess
from duckduckgo_search import DDGS

def web_search(query: str) -> str:
    """
    Search the web for information using DuckDuckGo.
    
    Args:
        query: The search query string.
        
    Returns:
        A string containing a summary of the search results.
    """
    try:
        results = DDGS().text(query, max_results=5)
        if not results:
            return "No results found."
        
        summary = ""
        for i, r in enumerate(results):
            summary += f"{i+1}. {r.get('title', '')}\n{r.get('body', '')}\nURL: {r.get('href', '')}\n\n"
        return summary
    except Exception as e:
        return f"Error performing web search: {str(e)}"

def read_file(file_path: str) -> str:
    """
    Read the contents of a file on the local system.
    
    Args:
        file_path: The absolute or relative path to the file.
        
    Returns:
        The text content of the file.
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            return f.read()
    except Exception as e:
        return f"Error reading file: {str(e)}"

def write_file(file_path: str, content: str) -> str:
    """
    Write content to a file on the local system. If the file exists, it will be overwritten.
    
    Args:
        file_path: The absolute or relative path to the file.
        content: The text content to write into the file.
        
    Returns:
        A success or error message.
    """
    try:
        # Create directories if they don't exist
        os.makedirs(os.path.dirname(os.path.abspath(file_path)), exist_ok=True)
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        return f"Successfully wrote to {file_path}"
    except Exception as e:
        return f"Error writing file: {str(e)}"

def execute_system_command(command: str) -> str:
    """
    Execute a shell/system command and return its output.
    
    Args:
        command: The command line string to execute.
        
    Returns:
        The standard output and standard error of the command.
    """
    try:
        result = subprocess.run(
            command, 
            shell=True, 
            capture_output=True, 
            text=True,
            timeout=30 # Prevent hanging
        )
        output = result.stdout
        if result.stderr:
            output += f"\nSTDERR:\n{result.stderr}"
        return output if output.strip() else "Command executed successfully with no output."
    except Exception as e:
        return f"Error executing command: {str(e)}"

def calculate(expression: str) -> str:
    """
    Evaluate a mathematical expression.
    
    Args:
        expression: A mathematical expression string (e.g. '123 * 456', '2**10').
        
    Returns:
        The result of the calculation.
    """
    try:
        # Use eval with a restricted environment for safety
        allowed_names = {"__builtins__": None}
        result = eval(expression, allowed_names)
        return str(result)
    except Exception as e:
        return f"Error evaluating expression: {str(e)}"

# List of all available tools for the agent
AVAILABLE_TOOLS = [
    web_search,
    read_file,
    write_file,
    execute_system_command,
    calculate
]
