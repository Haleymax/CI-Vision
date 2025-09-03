from langchain.agents import Tool

from mcp.file_opt import list_directory

# 定义工具列表
file_tools = [
    Tool(
        name="list_directory",
        description="List files in a directory",
        func=lambda path: list_directory(path, recursive=False)
    ),
    Tool(
        name="list_files_recursive", 
        description="List files in a directory recursively",
        func=lambda path: list_directory(path, recursive=True)
    )
]