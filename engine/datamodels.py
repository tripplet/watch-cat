from datetime import datetime
from google.appengine.ext import db

class Action(db.Model):
  enabled        = db.BooleanProperty(required=True, default=False)
  last_performed = db.DateTimeProperty()
  message        = db.StringProperty()


  def updatePerformed(self):
    self.last_performed = datetime.now()
    self.put()


  def performAction(self):
    pass


class PollEvent(db.Model):
  host       = db.StringProperty(required=True)
  port       = db.IntegerProperty(required=True)
  poll_type  = db.StringProperty(default='http', choices=set(['raw', 'http']))


class EmailAction(Action, db.Model):
  address    = db.EmailProperty(required=True)


class WebrequestAction(Action):
  url        = db.LinkProperty(required=True)


class IMAction(Action):
  im         = db.IMProperty(required=True)

