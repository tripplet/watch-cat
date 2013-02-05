#!/usr/bin/env python

import webapp2

from BlockedIP import BlockedIP
from WatchJob import WatchJob

class JobHandler(webapp2.RequestHandler):
  def get(self, key):
    if BlockedIP.isRemoteBlocked(self.request.remote_addr):
      self.abort(400)

    job = WatchJob.all().filter('secret =', key).get()

    # invalid key
    if job is None:
      BlockedIP.updateBlocked(self.request.remote_addr)
      self.abort(400)

    job.update(self.request.remote_addr)


# Main handler ######################################################################
# ###################################################################################
app = webapp2.WSGIApplication([('/job/(\w*)', JobHandler)], debug=False)
# ###################################################################################
# ###################################################################################