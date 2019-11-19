from datetime import timedelta, datetime
from google.appengine.ext import db


class LogEntry(db.Model):
    job = db.ReferenceProperty()
    name = db.StringProperty(required=True)
    event_time = db.DateTimeProperty(required=True, auto_now_add=True)
    detail = db.StringProperty()

    @staticmethod
    def log_event(job, name, detail):
        entry = LogEntry(job=job, name=name, detail=detail)
        entry.put()

    @staticmethod
    def cleanup():
        to_delete = [] 
        for evt in LogEntry.all(keys_only=True).filter('event_time < ', datetime.utcnow() - timedelta(days=365)):
            to_delete.append(evt)
        
        if len(to_delete) > 0:
            db.delete(to_delete)
