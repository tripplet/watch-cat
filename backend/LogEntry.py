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
