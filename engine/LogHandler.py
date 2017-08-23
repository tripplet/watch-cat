#!/usr/bin/env python

import webapp2
import jinja2
import os
from LogEntry import LogEntry
from WatchJob import WatchJob
from German_tzinfo import German_tzinfo

jinja_environment = jinja2.Environment(loader=jinja2.FileSystemLoader(os.path.dirname(__file__)))

class LogHandler(webapp2.RequestHandler):
    @staticmethod
    def formatDateTime(value):
      german_time = German_tzinfo()
      if value == None:
        return 'Never'
      else:
        return german_time.utcToLocal(value).strftime('%d.%m.%Y - %H:%M:%S')

    def get(self, job_name):
      template = jinja_environment.get_template('templates/log_template.htm')

      job_id = WatchJob.all().filter('name =', job_name).get()
      job_logs = LogEntry.all().filter('job =', job_id).order('-event_time').run()

      self.response.out.write(template.render(name = job_name, logging=job_logs))

jinja_environment.filters['formatDateTime'] = LogHandler.formatDateTime

app = webapp2.WSGIApplication([('/log/(\w+)', LogHandler)], debug=False)
