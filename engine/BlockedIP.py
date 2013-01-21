from datetime import datetime
from google.appengine.ext import db

class BlockedIP(db.Model):
  remote_ip        = db.StringProperty(required=True)
  invalid_requests = db.IntegerProperty(default=0)
  last_invalid     = db.DateTimeProperty()
  blocked_until    = db.DateTimeProperty()

  @staticmethod
  def isRemoteBlocked(remote_ip):
    entry = BlockedIP.all().filter('remote_ip =', remote_ip).get()

    if entry is None or entry.blocked_until is None:
      return False
    else:
      return (datetime.now() < entry.blocked_until)

  @staticmethod
  def updateRemote(remote_ip):
    entry = BlockedIP.all().filter('remote_ip =', remote_ip).get()

    if entry is None:
      entry = BlockedIP(remote_ip=remote_ip, invalid_requests=1)
    else:
      entry.invalid_requests = entry.invalid_requests + 1

    entry.last_invalid = datetime.now()
    entry.put()

