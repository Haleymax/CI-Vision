import yaml

from config.path import config_file

app_config = None

class AppConfig:
    _instance = None

    def __new__(cls, config_file="dev_config.yaml"):
        if cls._instance is None:
            cls._instance = super(AppConfig, cls).__new__(cls)
            cls._instance.config_file = config_file
            cls._instance._load_config()
        return cls._instance

    def _load_config(self):
        """加载配置文件"""
        try:
            with open(self.config_file, 'r') as f:
                self.config_data = yaml.safe_load(f) or {}
        except (FileNotFoundError, yaml.YAMLError):
            self.config_data = {}

    def get(self, key, default=None):
        """获取配置值"""
        return self.config_data.get(key, default)

    def __getitem__(self, key):
        """字典式访问"""
        return self.config_data[key]

    def reload(self):
        """重新加载配置"""
        self._load_config()

def get_config():
    """获取全局配置实例"""
    global app_config
    if app_config is None:
        app_config = AppConfig(config_file)
    return app_config