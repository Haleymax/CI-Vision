import logging
from django.utils import timezone

from django.db import models
from django.db.models import Count, Q

# Create your models here.

logger = logging.getLogger('django')

class JenkinsJob(models.Model):
    name = models.CharField(max_length=255, verbose_name='Job Name', unique=True)
    description = models.TextField(blank=True, verbose_name='Job Description')
    label = models.CharField(max_length=255, verbose_name='Job Label')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name='Created At')
    updated_at = models.DateTimeField(auto_now=True, verbose_name='Updated At')

    class Meta:
        verbose_name = 'Jenkins Job'
        verbose_name_plural = 'Jenkins Jobs'
        ordering = ['-created_at']

    def __str__(self):
        return self.name

    def latest_build(self):
        """
        Get the latest build for this job.
        """
        return self.logs.order_by('-build_number').first()

    def get_build_by_number(self, build_number):
        """
        Get a specific build by its number.
        """
        return self.logs.filter(build_number=build_number).first()

class JenkinsBuild(models.Model):
    BUILD_STATUS_CHOICES = (
        ("SUCCESS", "成功"),
        ("FAILURE", "失败"),
        ("RUNNING", "运行中"),
        ("ABORTED", "已中止"),
        ("UNSTABLE", "不稳定"),
    )
    job = models.ForeignKey(JenkinsJob, on_delete=models.CASCADE, related_name='build')
    status = models.CharField(max_length=20, choices=BUILD_STATUS_CHOICES, verbose_name='Build Status', null=True, blank=True)
    build_number = models.CharField(max_length=50)
    log_url = models.URLField(max_length=500, null=True, blank=True)
    start_time = models.DateTimeField(verbose_name='Start Time', null=True, blank=True)
    duration = models.DurationField(verbose_name='Duration', null=True, blank=True)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name='Created At')
    updated_at = models.DateTimeField(auto_now=True, verbose_name='Updated At')

    class Meta:
        unique_together = ('job', 'build_number')
        ordering = ['-start_time']

    def __str__(self):
        return f"{self.job.name} - Build {self.build_number} ({self.status})"


    @property
    def is_successful(self):
        return self.status == "SUCCESS"

class CompileBuildQuerySet(models.QuerySet):
    def by_branch(self, branch):
        return self.filter(branch=branch)

    def by_build_type(self, build_type):
        return self.filter(build_type=build_type)

    def successful(self):
        return self.filter(compile_status="SUCCESS")

    def failed(self):
        return self.filter(compile_status="FAILURE")

    def recent_days(self, days=7):
        from datetime import timedelta
        since = timezone.now() - timedelta(days=days)
        return self.filter(build_time__gte=since)

    def with_test_summary(self):
        return self.prefetch_related('test_builds').annotate(
            total_tests=Count('test_builds', filter=Q(test_builds__is_latest=True)),
            passed_tests= Count('test_builds', filter=Q(
                test_builds__is_latest=True,
                test_builds__test_result='SUCCESS'
            )),
            failed_tests= Count('test_builds', filter=Q(
                test_builds__is_latest=True,
                test_builds__test_result='FAILURE'
            ))
        )

class CompileBuild(models.Model):
    BUILD_TYPE_CHOICES = (
        ("release", "Release"),
        ("debug", "Debug"),
        ("kasan", "Kasan"),
        ("swkasan", "Swkasan"),
        ("asan", "Asan"),
    )

    COMPILE_STATUS_CHOICES = (
        ("SUCCESS", "Success"),
        ("FAILURE", "Failure"),
    )

    jenkins_build = models.ForeignKey(
        JenkinsBuild,
        on_delete=models.CASCADE,
        related_name='compile_builds',
    )
    version = models.CharField(max_length=50, verbose_name='Version')
    build_type = models.CharField(
        max_length=20,
        choices=BUILD_TYPE_CHOICES,
        default='release',
    )
    branch = models.CharField(max_length=100, verbose_name='Branch')
    compile_status = models.CharField(
        max_length=20,
        choices=COMPILE_STATUS_CHOICES,
        default='SUCCESS',
    )
    duration = models.DurationField(null=True)
    artifact_url = models.URLField(max_length=500)
    artifact_name = models.CharField(max_length=255, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    objects = CompileBuildQuerySet.as_manager()

    def save(self, *args, **kwargs):
        if self.artifact_url and not self.artifact_name:
            self.artifact_name = self.artifact_url.split('/')[-1]
        super().save(*args, **kwargs)

    class Meta:
        ordering = ['-created_at']

    def __str__(self):
        return f"{self.version} - {self.build_type} - {self.artifact_name}"

    @property
    def is_successful(self):
        return self.compile_status == "SUCCESS"


class TestBuildQuerySet(models.QuerySet):
    """测试构建查询集"""

    def latest_only(self):
        return self.filter(is_latest=True)

    def by_test_type(self, test_type):
        return self.filter(test_type=test_type)

    def successful(self):
        return self.filter(test_result="SUCCESS")

    def failed(self):
        return self.filter(test_result="FAILURE")


class TestBuild(models.Model):

    TEST_RESULT_CHOICES = (
        ('SUCCESS', 'Success'),
        ('FAILURE', 'Failure'),
        ('UNSTABLE', 'Unstable'),
        ('ABORTED', 'Aborted'),
    )

    compile_build = models.ForeignKey(CompileBuild, on_delete=models.CASCADE, related_name='test_builds')
    jenkins_build = models.ForeignKey(JenkinsBuild, on_delete=models.PROTECT, related_name='test_builds')

    test_type = models.CharField(max_length=100, null=True, blank=True)
    test_result = models.CharField(max_length=20, choices=TEST_RESULT_CHOICES)
    total_cases = models.PositiveIntegerField(default=0)
    passed_count = models.PositiveIntegerField(default=0)
    failed_count = models.PositiveIntegerField(default=0)

    failed_cases = models.JSONField(null=True, blank=True)
    default_ids = models.JSONField(null=True, blank=True)

