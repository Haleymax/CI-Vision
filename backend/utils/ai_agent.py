from langchain_deepseek import ChatDeepSeek
from langchain.agents import Tool, initialize_agent, AgentType
from mcp.tools import file_tools

class AIAgent:
    def __init__(self, model, api_key):
        self.llm = ChatDeepSeek(
            model=model,
            temperature=0.7,
            max_tokens=2048,
            api_key=api_key
        )


class FileAgent(AIAgent):
    def __init__(self, model, api_key):
        super().__init__(model, api_key)
        self.agent = initialize_agent(
            tools = file_tools,
            llm=self.llm,
            agent=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
            verbose=True
        )