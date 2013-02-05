#!/usr/bin/env python

import webapp2
from WatchJob import WatchJob
from PushOverAction import PushOverAction
from EmailAction import EmailAction
from datetime import datetime


class MainHandler(webapp2.RequestHandler):
    def get(self):
      jobs = WatchJob.all()
      self.response.write('<b><i>ServerTime: </i></b>' + datetime.now().strftime('%H:%M:%S') + '<br><br>')

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

      action1 = PushOverAction(enabled=True, token='***REMOVED***', message='test')
      action1.userkeys = ['***REMOVED***']
      action1.put()

      action2 = EmailAction(enabled=True, address = '***REMOVED***', subject='watch-cat: ', message='test')
      action2.put()

      new_job.actions = [action1.key(), action2.key()]
      new_job.generateSecret()
      new_job.put()

      self.response.write('Done')


class NofifyTest(webapp2.RequestHandler):
    def get(self, job_name):
      WatchJob.testJobActions(job_name)


app = webapp2.WSGIApplication([('/', MainHandler),
                               ('/create', CreateJob),
                               ('/notify/(\w+)', NofifyTest)], debug=False)
