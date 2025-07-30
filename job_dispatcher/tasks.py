import logging
from celery import shared_task

from job_dispatcher.models import JenkinsTask, JenkinsBuild
from utils.jenkins import JenkinsClient

logger = logging.getLogger('celery')

@shared_task
def trigger_jenkins_job(task_id, trigger_type, build_id=None):
    logger.info(f"Triggering jenkins job start for{task_id}")
    task = JenkinsTask.objects.get(id=task_id)
    if task is None:
        logger.error(f"Task with id {task_id} not found.")
        return
    try:
        client = JenkinsClient()
        client.set_logger(logger)
        if trigger_type == "jenkins":
            build_info = client.build_info(task.job_name, build_id)
        else:
            build_info = client.trigger_job(task.job_name, task.parameters)
        JenkinsBuild.objects.create(
            task=task,
            job_name= task.job_name,
            url=build_info.get_url(),
            number=build_info.get_number(),
            stages="RUNNING" if build_info.is_running() else build_info.get_stages(),
            parameters=build_info.get_parameters(),
        )
        logger.info(f"Trigger Jenkins job END for task: {build_info}")
    except Exception as e:
        logger.error(f"Error triggering Jenkins job: {e}")
        JenkinsBuild.objects.create(
            task=task,
            job_name= task.job_name,
            url="",
            number=-1,
            parameters={"ERROR": str(e)},
            status="ERROR"
        )