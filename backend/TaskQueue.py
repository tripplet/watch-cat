import webapp2
from google.appengine.ext import db


# TaskHandler
class TaskHandler(webapp2.RequestHandler):
    def post(self):
        key = self.request.get('key')
        db.get(key).check()


# Web handler
app = webapp2.WSGIApplication([('/task', TaskHandler)], debug=False)
