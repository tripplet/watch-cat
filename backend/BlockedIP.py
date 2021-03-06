from datetime import timedelta, datetime
from google.appengine.ext import db


class BlockedIP(db.Model):
    remote_ip = db.StringProperty(required=True)
    invalid_requests = db.IntegerProperty(default=0)
    last_invalid = db.DateTimeProperty()
    blocked_until = db.DateTimeProperty()

    @staticmethod
    def is_remote_blocked(remote_ip):
        entry = BlockedIP.all().filter('remote_ip =', remote_ip).get()

        if entry is None or entry.blocked_until is None:
            return False
        else:
            return datetime.utcnow() < entry.blocked_until

    @staticmethod
    def update_blocked(remote_ip):
        entry = BlockedIP.all().filter('remote_ip =', remote_ip).get()

        if entry is None:
            entry = BlockedIP(remote_ip=remote_ip, invalid_requests=1)
        else:
            entry.invalid_requests = entry.invalid_requests + 1

            # block ip after 10 invalid requests for 1 hour
            if entry.invalid_requests >= 10:
                entry.blocked_until = datetime.utcnow() + timedelta(minutes=1)

        entry.last_invalid = datetime.utcnow()
        entry.put()

    @staticmethod
    def remove_outdated():
        for p in BlockedIP.all().filter('blocked_until < ', datetime.utcnow()):
            p.delete()
