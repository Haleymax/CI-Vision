from rest_framework import serializers

from job_dispatcher.models import JenkinsJob, JenkinsBuild, JenkinsTask


class JenkinsJobSerializer(serializers.ModelSerializer):
    class Meta:
        model = JenkinsJob
        fields = '__all__'

class JenkinsTaskSerializer(serializers.ModelSerializer):
    class Meta:
        model = JenkinsTask
        fields = '__all__'

class JenkinsBuildSerializer(serializers.ModelSerializer):
    class Meta:
        model = JenkinsBuild
        fields = '__all__'