#!/usr/bin/env python

import webapp2
import jinja2
import os
from WatchJob import WatchJob
from PushOverAction import PushOverAction
from EmailAction import EmailAction
from datetime import datetime

jinja_environment = jinja2.Environment(loader=jinja2.FileSystemLoader(os.path.dirname(__file__)))

class MobileHandler(webapp2.RequestHandler):
    @staticmethod
    def formatDateTime(value):
      if value == None:
        return 'Never'
      else:
        return value.strftime('%H:%M:%S - %d.%m.%Y')

    def get(self):
      template = jinja_environment.get_template('templates/main_template.htm')

      jobs = WatchJob.all().run()

      template_values = {
        'timestring': datetime.now().strftime('%H:%M:%S'),
        'jobs': jobs
      }

      self.response.out.write(template.render(template_values))

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

jinja_environment.filters['formatDateTime'] = MobileHandler.formatDateTime

app = webapp2.WSGIApplication([('/', MainHandler),
                               ('/m', MobileHandler),
                               ('/create', CreateJob),
                               ('/notify/(\w+)', NofifyTest)], debug=False)
