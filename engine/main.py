#!/usr/bin/env python

import webapp2
from WatchJob import WatchJob
from PushOverAction import PushOverAction
from datetime import datetime


class MainHandler(webapp2.RequestHandler):
    def get(self):
      jobs = WatchJob.all()
      self.response.write('<b><i>ServerTime: </i></b>' + datetime.now().strftime('%H:%M:%S - %d.%m.%Y') + '<br><br>')

      for job in jobs:
        if job.last_seen == None:
          job_lastseen = 'Never'
        else:
          job_lastseen = job.last_seen.strftime('%H:%M:%S - %d.%m.%Y')

        self.response.write('<b>%s</b> <a href="/notify/%s">[testNotification]</a><br>%s<br>%s<br><br>' %
          (job.name, job.name, job_lastseen, job.last_ip))


class CreateJob(webapp2.RequestHandler):
    def get(self):
      new_job = WatchJob(name='new_job', interval=30)

      action = PushOverAction(enabled=True, token='', message='test')
      action.userkeys = ['']
      action.put()

      new_job.actions = [action.key()]
      new_job.generateSecret()
      new_job.put()

      self.response.write('Done')


class NofifyTest(webapp2.RequestHandler):
    def get(self, job_name):
      WatchJob.testJobActions(job_name)


app = webapp2.WSGIApplication([('/', MainHandler),
                               ('/create', CreateJob),
                               ('/notify/(\w+)', NofifyTest)], debug=False)
