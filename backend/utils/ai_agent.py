import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(__file__)))

from langchain_deepseek import ChatDeepSeek
from langchain.agents import initialize_agent, AgentType
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
            tools=file_tools,
            llm=self.llm,
            agent=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
            verbose=True
        )

    def list_all_files(self, file_path):
        """
        列出目录中的所有文件并使用AI添加解释
        """
        try:
        # 改进的提示词，更清晰地指导AI使用工具
            prompt = f"""
            请帮我分析目录 {file_path} 中的文件：
            1. 首先使用工具列出该目录下的所有文件
            2. 然后根据文件扩展名分析每个文件的类型和用途
            3. 最后提供一个简洁的总结

            请使用可用的工具来获取文件列表。
            """
        
            result = self.agent.run(prompt)
            return result
        
        except Exception as e:
            return f"分析文件时出错: {str(e)}"