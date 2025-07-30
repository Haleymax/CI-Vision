from django.urls import path, include
from rest_framework.routers import DefaultRouter

from job_dispatcher.views import JenkinsJobViewSet, JenkinsTaskViewSet

router = DefaultRouter()
router.register(f"jobs", JenkinsJobViewSet, basename="jobs")
router.register(r"tasks", JenkinsTaskViewSet, basename="tasks")

urlpatterns = [
    path("", include(router.urls)),
]