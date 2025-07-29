import logging

from django.contrib.admin import action
from rest_framework import viewsets
from rest_framework.permissions import AllowAny

from job_dispatcher.models import JenkinsJob, JenkinsTask
from job_dispatcher.serializers import JenkinsJobSerializer, JenkinsBuildSerializer

# Create your views here.


logger = logging.getLogger('django')

class JenkinsJobViewSet(viewsets.ModelViewSet):
    queryset = JenkinsJob.objects.all()
    serializer_class = JenkinsJobSerializer
    permission_classes = [AllowAny]


class JenkinsTaskViewSet(viewsets.ModelViewSet):
    """
    A viewset for viewing and editing Jenkins tasks.
    """
    queryset = JenkinsTask.objects.prefetch_related('builds').order_by('id').all()
    serializer_class = JenkinsBuildSerializer
    permission_classes = [AllowAny]
    ordering_fields = ['id', 'created_at']
    ordering = ['-id']

    def get_queryset(self):
        qs = super().get_queryset()
        title = self.request.query_params.get('title', None)
        user = self.request.query_params.get('user', None)
        mine = self.request.query_params.get('mine', None)

        if title:
            qs = qs.filter(title__icontains=title)
        if user:
            qs = qs.filter(user__username__icontains=user)
        if mine:
            qs = qs.filter(user=self.request.user.username)

        return qs

    @action(detail=False, methods=['post'])
    def trigger_job(self, request, *args, **kwargs):
        """
        Trigger a Jenkins job.
        """
        job_name = request.data.get('job_name')
        params = request.data.get('params', {})

        if not job_name:
            return Response({"error": "Job name is required."}, status=400)

        try:
            jenkins_client = JenkinsClient()
            build = jenkins_client.trigger_job(job_name, params=params)
            return Response({"message": "Job triggered successfully.", "build_id": build.id}, status=200)
        except Exception as e:
            logger.error(f"Failed to trigger job {job_name}: {str(e)}")
            return Response({"error": str(e)}, status=500)