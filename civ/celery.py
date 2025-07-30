import os

from celery import Celery
from celery.signals import setup_logging

# 设置 Django 的默认配置模块
os.environ.setdefault("DJANGO_SETTINGS_MODULE", "civ.settings")

app = Celery("civ")

#　从 Django 的配置文件中加载 Celery 配置
app.config_from_object("django.conf:settings", namespace="CELERY")

@setup_logging.connect
def setup_logging(*args, **kwargs):
    """禁用 Celery 的日志配置，使用 Django 的日志配置"""
    from logging.config import dictConfig
    from django.conf import settings
    dictConfig(settings.LOGGING)

# 自动发现所有 Django 应用中的任务
app.autodiscover_tasks()