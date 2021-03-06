#!/usr/bin/env python

import webapp2

from BlockedIP import BlockedIP
from WatchJob import WatchJob


class JobHandler(webapp2.RequestHandler):
    def get(self):
        if BlockedIP.is_remote_blocked(self.request.remote_addr):
            self.abort(400)
            return

        try:
            key = self.request.get('key')
        except:
            self.abort(400)
            return

        try:
            uptime = int(self.request.get('uptime'))
        except:
            uptime = None

        job = WatchJob.all().filter('secret =', key).get()

        # invalid key
        if job is None:
            BlockedIP.update_blocked(self.request.remote_addr)
            self.abort(400)
            return

        job.update(self.request.remote_addr, uptime)


# Main handler
app = webapp2.WSGIApplication([('/job', JobHandler)], debug=False)
