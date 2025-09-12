from langchain_deepseek import ChatDeepSeek
from langchain_core.messages import HumanMessage, SystemMessage

from config.config import get_config

app_config = get_config().get("llm", {})


def get_deepseek_llm():
    return ChatDeepSeek(
        model=app_config.get("model", "deepseek-chat"),
        temperature=app_config.get("temperature", 0.7),
        max_tokens=app_config.get("max_tokens", 2048),
        api_key=app_config.get("api_key")
    )


llms = {
    "deepseek": get_deepseek_llm()
}


def get_llm(model_name: str):
    return llms.get(model_name)