#!/usr/bin/env python

import webapp2
import jinja2
import os
from WatchJob import WatchJob
from GermanTzInfo import GermanTzInfo
from PushOverAction import PushOverAction
from EmailAction import EmailAction
from datetime import datetime, timedelta

jinja_environment = jinja2.Environment(loader=jinja2.FileSystemLoader(os.path.dirname(__file__)))


def format_datetime(value):
    german_time = GermanTzInfo()
    if value is None:
        return 'Never'
    else:
        return german_time.utc_to_local(value).strftime('%H:%M:%S - %d.%m.%Y')


def format_timespan(value):
    if value is None:
        return 'Unknown'
    else:
        return str(timedelta(seconds=value))

def add_breakchars(value):
    if value is None:
        return '&nbsp;'
    elif ':' in value:
        return value.replace(':', ':<wbr/>')
    elif '.' in value:
        return value.replace('.', '.<wbr/>')


class MainHandler(webapp2.RequestHandler):
    def get(self):
        german_time = GermanTzInfo()
        template = jinja_environment.get_template('templates/main_template.htm')

        jobs = WatchJob.all().run()

        template_values = {
            'timestring': german_time.utc_to_local(datetime.utcnow()).strftime('%H:%M:%S'),
            'jobs': jobs
        }

        self.response.out.write(template.render(template_values))


class DebugHandler(webapp2.RequestHandler):
    def get(self):
        german_time = GermanTzInfo()

        jobs = WatchJob.all()
        self.response.write('<b><i>ServerTime: </i></b>' + german_time.utc_to_local(datetime.utcnow()).strftime(
            '%H:%M:%S') + '<br><br>')

        for job in jobs:
            self.response.write('<b>%s</b> <a href="/notify/%s">[testNotification]</a><br>%s<br>%s<br>%s<br><br>' %
                                (job.name, job.name, format_datetime(job.last_seen), job.last_ip,
                                 format_timespan(job.uptime)))


class CreateJob(webapp2.RequestHandler):
    def get(self):
        new_job = WatchJob(name='new_job', interval=30)

        action1 = PushOverAction(enabled=True, token='***REMOVED***', message='test')
        action1.userkeys = ['***REMOVED***']
        action1.put()

        action2 = EmailAction(enabled=True, address='***REMOVED***', subject='watch-cat: ', message='test')
        action2.put()

        action3 = PushOverAction(enabled=True, token='***REMOVED***', message='test2')
        action3.userkeys = ['***REMOVED***']
        action3.put()

        new_job.timeout_actions = [action1.key(), action2.key()]
        new_job.backonline_actions = [action3.key()]
        new_job.generate_secret()
        new_job.put()

        self.response.write('Done')


class NofifyTest(webapp2.RequestHandler):
    def get(self, job_name):
        WatchJob.test_job_actions(job_name)


jinja_environment.filters['format_datetime'] = format_datetime
jinja_environment.filters['format_timespan'] = format_timespan
jinja_environment.filters['add_breakchars'] = add_breakchars

app = webapp2.WSGIApplication([('/', MainHandler),
                               ('/debug', DebugHandler),
                               ('/create', CreateJob),
                               ('/notify/(\w+)', NofifyTest)], debug=False)
