#!/usr/bin/env python

import webapp2
from WatchJob import WatchJob

class MainHandler(webapp2.RequestHandler):
    def get(self):
      jobs = WatchJob.all()

      for job in jobs:
        self.response.write(job.name)
        self.response.write('<br>')
        self.response.write(job.last_seen)
        self.response.write('<br>')
        self.response.write(job.last_ip)
        self.response.write('<br>')
        self.response.write('<br>')

class CreateJob(webapp2.RequestHandler):
    def get(self):
      new_job = WatchJob(name='new_job', interval=30)
      new_job.generateSecret()
      new_job.put()

      self.response.write('Done')

app = webapp2.WSGIApplication([('/', MainHandler),
                               ('/create', CreateJob)], debug=True)
