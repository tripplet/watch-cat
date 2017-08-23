#!/usr/bin/env python

import webapp2
import jinja2
import os
from WatchJob import WatchJob
from German_tzinfo import German_tzinfo
from PushOverAction import PushOverAction
from EmailAction import EmailAction
from datetime import datetime, timedelta

jinja_environment = jinja2.Environment(loader=jinja2.FileSystemLoader(os.path.dirname(__file__)))


def formatDateTime(value):
  german_time = German_tzinfo()
  if value == None:
    return 'Never'
  else:
    return german_time.utcToLocal(value).strftime('%H:%M:%S - %d.%m.%Y')


def formatTimespan(value):
  if value == None:
    return 'Unknown'
  else:
    return str(timedelta(seconds=value))


class MainHandler(webapp2.RequestHandler):
    def get(self):
      german_time = German_tzinfo()
      template = jinja_environment.get_template('templates/main_template.htm')

      jobs = WatchJob.all().run()

      template_values = {
        'timestring': german_time.utcToLocal(datetime.utcnow()).strftime('%H:%M:%S'),
        'jobs': jobs
      }

      self.response.out.write(template.render(template_values))

class DebugHandler(webapp2.RequestHandler):
    def get(self):
      german_time = German_tzinfo()

      jobs = WatchJob.all()
      self.response.write('<b><i>ServerTime: </i></b>' + german_time.utcToLocal(datetime.utcnow()).strftime('%H:%M:%S') + '<br><br>')

      for job in jobs:
        self.response.write('<b>%s</b> <a href="/notify/%s">[testNotification]</a><br>%s<br>%s<br>%s<br><br>' %
          (job.name, job.name, formatDateTime(job.last_seen), job.last_ip, formatTimespan(job.uptime)))


class CreateJob(webapp2.RequestHandler):
    def get(self):
      new_job = WatchJob(name='new_job', interval=30)

      action1 = PushOverAction(enabled=True, token='***REMOVED***', message='test')
      action1.userkeys = ['***REMOVED***']
      action1.put()

      action2 = EmailAction(enabled=True, address = '***REMOVED***', subject='watch-cat: ', message='test')
      action2.put()

      action3 = PushOverAction(enabled=True, token='***REMOVED***', message='test2')
      action3.userkeys = ['***REMOVED***']
      action3.put()

      new_job.timeout_actions = [action1.key(), action2.key()]
      new_job.backonline_actions = [action3.key()]
      new_job.generateSecret()
      new_job.put()

      self.response.write('Done')


class NofifyTest(webapp2.RequestHandler):
    def get(self, job_name):
      WatchJob.testJobActions(job_name)

jinja_environment.filters['formatDateTime'] = formatDateTime
jinja_environment.filters['formatTimespan'] = formatTimespan


app = webapp2.WSGIApplication([('/', MainHandler),
                               ('/debug', DebugHandler),
                               ('/create', CreateJob),
                               ('/notify/(\w+)', NofifyTest)], debug=False)
