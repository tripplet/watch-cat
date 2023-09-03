#!/usr/bin/env python

import webapp2
import jinja2
import os
from LogEntry import LogEntry
from WatchJob import WatchJob
from GermanTzInfo import GermanTzInfo

jinja_environment = jinja2.Environment(loader=jinja2.FileSystemLoader(os.path.dirname(__file__)))


class LogHandler(webapp2.RequestHandler):
    @staticmethod
    def format_datetime(value):
        german_time = GermanTzInfo()
        if value is None:
            return 'Never'
        else:
            return german_time.utc_to_local(value).strftime('%d.%m.%Y - %H:%M:%S')

    @staticmethod
    def add_breakchars(value):
        return value.replace(':', ':<wbr/>').replace('.', '.<wbr/>')

    def get(self, job_name):
        template = jinja_environment.get_template('templates/log_template.htm')

        job_id = WatchJob.all().filter('name =', job_name).get()
        job_logs = LogEntry.all().filter('job =', job_id).order('-event_time').run(limit=100)

        self.response.out.write(template.render(name=job_name, logging=job_logs))


jinja_environment.filters['format_datetime'] = LogHandler.format_datetime
jinja_environment.filters['add_breakchars'] = LogHandler.add_breakchars

app = webapp2.WSGIApplication([
    webapp2.Route(r'/log/<job_name>', handler=LogHandler, name="log-handler"),
], debug=False)
