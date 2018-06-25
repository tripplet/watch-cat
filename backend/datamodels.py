from datetime import datetime
from google.appengine.ext import db


class Action(db.Model):
    enabled = db.BooleanProperty(required=True, default=False)
    last_performed = db.DateTimeProperty()
    failure_action = db.BooleanProperty(default=True)
    message = db.StringProperty()

    def update_performed(self):
        self.last_performed = datetime.utcnow()
        self.put()

    def perform_action(self):
        pass


class PollEvent(db.Model):
    host = db.StringProperty(required=True)
    port = db.IntegerProperty(required=True)
    poll_type = db.StringProperty(default='http', choices={'raw', 'http'})


class WebrequestAction(Action):
    url = db.LinkProperty(required=True)


class IMAction(Action):
    im = db.IMProperty(required=True)
