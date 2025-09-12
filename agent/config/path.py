import os

from dotenv import load_dotenv

load_dotenv()

env = os.getenv('ENV')
base_path = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
config_file = os.path.join(base_path, "config", f"{env}_config.yaml" if env else "dev_config.yaml")