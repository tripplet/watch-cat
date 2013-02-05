import os, hashlib
from datetime import datetime, timedelta

from google.appengine.ext import db


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


  def generateSecret(self):
    self.secret = hashlib.sha1(os.urandom(1024)).hexdigest()


  def update(self, remote_ip):
    self.last_seen = datetime.now()
    self.last_ip   = remote_ip
    self.put()


  @staticmethod
  def checkJobs():
    jobs = WatchJob.all().filter('enabled =', True).run()

    # check all enabled jobs
    for entry in jobs:

      # check if job is overdue
      if entry.last_seen + timedelta(minutes=entry.interval) < datetime.now():

        # perform all actions
        for action_key in entry.actions:
          db.get(action_key).performAction()


  @staticmethod
  def testJobActions(job_name):
    job = WatchJob.all().filter('name =', job_name).get()

    if job == None:
      return False

    for action_key in job.actions:
      db.get(action_key).performAction()

    return True
