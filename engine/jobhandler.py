#!/usr/bin/env python

import webapp2
from datetime import datetime

from BlockedIP import BlockedIP
from WatchJob import WatchJob

class JobHandler(webapp2.RequestHandler):
  def get(self, key):
    if BlockedIP.isRemoteBlocked(self.request.remote_addr):
      self.abort(400)

    job = WatchJob.all().filter('secret =', key).get()

    # invalid key
    if job is None :
      BlockedIP.setBlocked(self.request.remote_addr)
      self.abort(400)

    job.last_seen = datetime.now()
    job.last_ip   = self.request.remote_addr

    job.put()


# Main handler ######################################################################
# ###################################################################################
app = webapp2.WSGIApplication([('/job/(\w*)', JobHandler)], debug=False)
# ###################################################################################
# ###################################################################################