import hashlib
import logging
import os
from datetime import datetime, timedelta

from google.appengine.api import taskqueue
from google.appengine.api.taskqueue import Queue, Task
from google.appengine.ext import db

# import all possible actions so .performAction() in check() works
from LogEntry import LogEntry
from PushOverAction import PushOverAction # DO NOT REMOVE
from EmailAction import EmailAction # DO NOT REMOVE


class WatchJob(db.Model):
    name = db.StringProperty(required=True)
    watch_type = db.StringProperty(required=True, default='push', choices={'push', 'poll'})
    enabled = db.BooleanProperty(required=True, default=False)
    interval = db.IntegerProperty(required=True)
    created = db.DateTimeProperty(required=True, auto_now_add=True)
    status = db.StringProperty(default='online', choices={'offline', 'online'})
    secret = db.StringProperty()
    last_fail = db.DateTimeProperty()
    last_seen = db.DateTimeProperty()
    last_ip = db.StringProperty()
    uptime = db.IntegerProperty()
    timeout_actions = db.ListProperty(db.Key)
    backonline_actions = db.ListProperty(db.Key)
    reboot_actions = db.ListProperty(db.Key)
    poll = db.ReferenceProperty()
    task_name = db.StringProperty()

    def generate_secret(self):
        self.secret = hashlib.sha1(os.urandom(1024)).hexdigest()

    def update(self, remote_ip, uptime):
        self.last_seen = datetime.utcnow()

        if self.last_ip != remote_ip:
            LogEntry.log_event(self.key(), 'Info', 'IP changed - new IP: ' + remote_ip)

        self.last_ip = remote_ip

        if uptime is not None:
            if self.update is not None and self.uptime > uptime:
                LogEntry.log_event(self.key(), 'Reboot',
                                   'Reboot - Previous uptime: ' + str(timedelta(seconds=self.uptime)))
                for action_key in self.reboot_actions:
                    try:
                        db.get(action_key).perform_action()
                    except Exception as exp:
                        logging.error('Error executing reboot action: ' + str(exp))

        self.uptime = uptime
        self.put()

        # job got back online
        if self.status == 'offline':
            self.status = 'online'
            LogEntry.log_event(self.key(), 'Info', 'Job back online - IP: ' + remote_ip)

            # perform all back_online actions
            for action_key in self.backonline_actions:
                try:
                    db.get(action_key).perform_action()
                except Exception as exp:
                    logging.error('Error executing backonline action: ' + str(exp))

        # delete previous (waiting) task
        if self.task_name is not None:
            logging.debug('old task: ' + self.task_name)
            Queue.delete_tasks(Queue(), Task(name=self.task_name))

        task_name = self.name + '_' + datetime.utcnow().strftime('%Y-%m-%d_%H-%M-%S-%f')

        # create task to be executed in updated no called in interval minutes
        taskqueue.add(name=task_name, url='/task', params={'key': self.key()}, countdown=(self.interval + 2) * 60)

        self.task_name = task_name
        self.put()

    def check(self):
        # check if job is overdue
        if self.last_seen + timedelta(minutes=self.interval) < datetime.utcnow():

            self.status = 'offline'
            LogEntry.log_event(self.key(), 'Error', 'job is overdue')

            # perform all actions
            for action_key in self.timeout_actions:
                try:
                    db.get(action_key).perform_action()
                except Exception as exp:
                    logging.error('Error executing timeout action: ' + str(exp))
            self.put()

    @staticmethod
    def check_all_jobs():
        jobs = WatchJob.all().filter('enabled =', True).run()

        # check all enabled jobs
        for entry in jobs:
            entry.check()

    @staticmethod
    def test_job_actions(job_name):
        job = WatchJob.all().filter('name =', job_name).get()

        if job is None:
            return False

        for action_key in job.timeout_actions:
            db.get(action_key).perform_action()

        return True
