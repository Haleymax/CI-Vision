import logging

from rest_framework.decorators import action
from rest_framework.response import Response
from rest_framework import viewsets, status
from rest_framework.permissions import AllowAny

from job_dispatcher.models import JenkinsJob, JenkinsTask
from job_dispatcher.serializers import JenkinsJobSerializer, JenkinsBuildSerializer
from job_dispatcher.tasks import trigger_jenkins_job

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
    def create_task(self, request, *args, **kwargs):
        """
        create a Jenkins job.
        """
        request_body = request.data
        title = request_body.get('title', None)
        job_name = request_body.get('job_name', None)
        parameters = request_body.get('parameters', {})
        if request.user and request.user.is_authenticated:
            user = request.user.username
        else:
            user = request_body.get('user', "default_user")

        build_id = request_body.get('build_id', None)
        trigger_type = request_body.get('trigger_type', None)

        result = JenkinsTask.objects.create(
            title=title,
            job_name=job_name,
            user=user,
            parameters=parameters,
            trigger_type=trigger_type
        )
        trigger_jenkins_job.delay(result.id, trigger_type, build_id)
        return Response(status=status.HTTP_200_OK, data={
            "id": result.id,
            "title": result.title,
        })