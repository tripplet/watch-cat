import os, hashlib
from datetime import datetime, timedelta
from google.appengine.ext import db

from datamodels import PollEvent, Action

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
  actions    = db.ReferenceProperty(Action)
  poll       = db.ReferenceProperty(PollEvent)

  def generateSecret(self):
    self.secret = hashlib.sha1(os.urandom(1024)).hexdigest()

  @staticmethod
  def checkJobs():
    jobs = WatchJob.all().filter('enabled =', True)

    for entry in jobs:
      if entry.last_seen + timedelta(minutes=entry.interval) < datetime.now():
        for action in entry.actions:
          action.performAction()