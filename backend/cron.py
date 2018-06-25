import webapp2
from BlockedIP import BlockedIP
from WatchJob import WatchJob


# Cronjob handler
class CronHandler(webapp2.RequestHandler):
    def get(self, action):
        if action == 'cleanup':
            BlockedIP.remove_outdated()
            self.response.out.write('Done!')
        elif action == 'checkjobs':
            WatchJob.check_all_jobs()
            self.response.out.write('Done!')


# Main handler
app = webapp2.WSGIApplication([('/cron/(\w+)', CronHandler)], debug=False)
