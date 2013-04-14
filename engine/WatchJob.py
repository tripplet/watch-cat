import os, hashlib
import logging
from datetime import datetime, timedelta

from google.appengine.api import taskqueue
from google.appengine.api.taskqueue import Queue, Task
from google.appengine.ext import db

# import all possible actions so .performAction() in check() works
from PushOverAction import PushOverAction
from EmailAction import EmailAction
from LogEntry import LogEntry

class WatchJob(db.Model):
  name       = db.StringProperty(required=True)
  watch_type = db.StringProperty(required=True, default='push', choices=set(['push', 'poll']))
  enabled    = db.BooleanProperty(required=True, default=False)
  interval   = db.IntegerProperty(required=True)
  created    = db.DateTimeProperty(required=True, auto_now_add=True)
  secret     = db.StringProperty()
  last_fail  = db.DateTimeProperty()
  last_seen  = db.DateTimeProperty()
  last_ip    = db.StringProperty()
  actions    = db.ListProperty(db.Key)
  poll       = db.ReferenceProperty()
  task_name  = db.StringProperty()


  def generateSecret(self):
    self.secret = hashlib.sha1(os.urandom(1024)).hexdigest()


  def update(self, remote_ip):
    self.last_seen = datetime.utcnow()
    self.last_ip   = remote_ip
    self.put()

    # delete previous (waiting) task
    if (self.task_name != None):
      logging.info('old task: ' + self.task_name)
      Queue.delete_tasks(Queue(), Task(name=self.task_name))

    task_name = self.name + '_' + datetime.utcnow().strftime('%Y-%m-%d_%H-%M-%S-%f')

    # create task to be executed in updated no called in interval minutes
    taskqueue.add(name=task_name, url='/task', params={'key': self.key()}, countdown=(self.interval + 1)*60)

    self.task_name = task_name
    self.put()


  def check(self):
    # check if job is overdue
    if self.last_seen + timedelta(minutes=self.interval) < datetime.utcnow():

      LogEntry.log_event(self.key() ,'Error', 'job is overdue')

      # perform all actions
      for action_key in self.actions:
        db.get(action_key).performAction()


  @staticmethod
  def checkAllJobs():
    jobs = WatchJob.all().filter('enabled =', True).run()

    # check all enabled jobs
    for entry in jobs:
      entry.check()


  @staticmethod
  def testJobActions(job_name):
    job = WatchJob.all().filter('name =', job_name).get()

    if job == None:
      return False

    for action_key in job.actions:
      db.get(action_key).performAction()

    return True
