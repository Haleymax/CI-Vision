from django.urls import path
from . import views

urlpatterns = [
    path('hello/', views.hello_world, name='hello_world'),
    path('demo/', views.api_demo, name='api_demo'),
    path('files/', views.list_files_with_ai, name='list_files_with_ai'),
]
