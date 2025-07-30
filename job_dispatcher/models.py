from django.db import models

# Create your models here.

def default_list():
    return []


class BaseModel(models.Model):
    """
    Base model for all models in the job dispatcher app.
    """
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        abstract = True # 不生成数据库表

class JenkinsTask(BaseModel):
    """
    Indicates the record that triggers the jenkins job.
    """
    TRIGGER_CHOICES = (
        ('civ', 'Manual'),
        ('api', 'api'),
        ('jenkins', 'Jenkins'),
    )
    title = models.CharField(max_length=200, default="", blank=True)
    job_name = models.CharField(max_length=200, default="", blank=True)
    user = models.CharField(max_length=200, default="", blank=True)
    parameters = models.JSONField(default=dict, blank=True)
    trigger_type = models.CharField(
        max_length=20,
        choices=TRIGGER_CHOICES,
        default='civ',
        blank=True
    )

    def __str__(self):
        return self.title


class JenkinsBuild(models.Model):
    STATUS_CHOICES = (
        ('SUCCESS', 'Success'),
        ('FAILURE', 'Failure'),
        ('UNSTABLE', 'Unstable'),
        ('ABORTED', 'Aborted'),
        ('NOT_BUILT', 'Not Built'),
        ('UNKNOWN', 'Unknown'),
    )
    task = models.ForeignKey(JenkinsTask, on_delete=models.CASCADE, related_name='builds')
    job_name = models.CharField(max_length=200, default="", blank=True)
    number = models.IntegerField(default=0)
    url = models.URLField(max_length=500, default="", blank=True)
    status = models.CharField(
        max_length=20,
        choices=STATUS_CHOICES,
        default='UNKNOWN',
        blank=True
    )
    start_time = models.DateTimeField(null=True, blank=True)
    end_time = models.DateTimeField(null=True, blank=True)
    duration = models.IntegerField(default=0)
    parameters = models.JSONField(default=dict, blank=True)
    stages = models.JSONField(default=default_list, blank=True)
    allure_report = models.TextField(default="", blank=True)
    artifact = models.JSONField(default=default_list, blank=True)


class JenkinsJob(models.Model):
    """
    Represents a Jenkins job.
    """
    name = models.CharField(max_length=200, unique=True)
    description = models.TextField(default="", blank=True)
    url = models.URLField(max_length=500, default="", blank=True)
    last_build_number = models.IntegerField(default=0)
    last_build_status = models.CharField(max_length=20, default="UNKNOWN", blank=True)

    def __str__(self):
        return self.name