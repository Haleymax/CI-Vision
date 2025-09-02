from langchain.agents import Tool

from mcp.file_opt import list_directory

file_tools = [
    Tool(
        name="list_directory",
        func=lambda path: list_directory(path, recursive=False)
    ),
    Tool(
        name="list_files_recursive",
        func=lambda path: list_directory(path, recursive=True)
    )
]