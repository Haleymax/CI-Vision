import logging
import threading
import time
import traceback
from urllib.parse import urljoin

from django.conf import settings
from jenkinsapi.custom_exceptions import NotBuiltYet
from jenkinsapi.jenkins import Jenkins
from requests import HTTPError

config = getattr(settings, 'JENKINS_CONFIG', {})


class JenkinsClient:
    _instance = None
    _lock = threading.Lock()
    _client = None

    def __init__(self, url=config.get('URL'),
                 username=config.get('USERNAME'),
                 password=config.get('password')):
        self.base_url = url
        if not self.base_url.endswith('/'):
            self.base_url = self.base_url + '/'
        self.username = username
        self.password = password
        self.logger = logging.getLogger('django')
        self._client = Jenkins(
            baseurl=self.base_url,
            username=self.username,
            password=self.password,
            timeout=60,
            lazy=True
        )
        self._init_auth()

    def _init_auth(self):
        pass

    def set_logger(self, logger):
        self.logger = logger

    def __new__(cls):
        with cls._lock:
            if not cls._instance:
                cls._instance = super().__new__(cls)
            return cls._instance

    def trigger_job(self, job_name, params=None, timeout=20, delay=2):
        job = self._client[job_name]
        job.get_queue_item
        build_queue = job.invoke(build_params=params)
        print('build_queue:', build_queue)
        time_taken = 0
        while time_taken < timeout:
            try:
                build_queue.poll()
                return build_queue.get_build()
            except NotBuiltYet:
                time.sleep(delay)
                time_taken += delay
                continue
            except HTTPError:
                self.logger.info(traceback.format_exc())
                time.sleep(delay)
                time_taken += delay
                continue
        else:
            raise Exception('Timeout while waiting for build to start')

    def get_pipeline_stages(self, job_name, build_number):
        api_url = f'blue/rest/organizations/jenkins/pipelines/{job_name}/runs/{build_number}/nodes/?limit=10000'
        try:
            stages_data = self.get(api_url).json()
            stages = []
            for stage in stages_data:
                links = stage.get('_links')
                self_link = links.get('self', {})
                if self_link:
                    self_url = self_link.get('href')
                    links['log'] = {'href': f'{self_url}log/'}
                stage_info = {
                    'id': stage.get('id'),
                    'displayName': stage.get('displayName'),
                    'state': stage.get('state'),
                    'result': stage.get('result'),
                    'durationInMillis': stage.get('durationInMillis'),
                    'startTime': stage.get('startTime') if stage.get('startTime') else None,
                    'type': stage.get('type'),  # STAGE
                    'edges': stage.get('edges'),  # 下一个节点 [{'id‘: 'xxx', 'type': 'xxx'}]
                    'firstParent': stage.get('firstParent'),
                    'links': links
                }
                stages.append(stage_info)
            return stages
        except Exception as e:
            self.logger.error(traceback.format_exc())
            raise e

    def get_stage_log(self, stage):
        try:
            log_url = stage.get('links', {}).get('log', {}).get('href', '')
            if not log_url:
                self.logger.error(f'No log URL found in stage {stage}')
                return ''
            stage_log = self.get(log_url).text
            return stage_log
        except Exception as e:
            self.logger.error(traceback.format_exc())
            raise e

    def get_log(self, log_url):
        try:
            return self.get(log_url).text
        except Exception as e:
            self.logger.error(traceback.format_exc())
            raise e

    def get(self, url):
        api_url = urljoin(self.base_url, url)
        return self._client.requester.get_and_confirm_status(api_url)

    def build_info(self, job_name, number):
        try:
            return self._client[job_name][number]
        except Exception as e:
            self.logger.error(traceback.format_exc())
            raise e