#!/usr/bin/env python

import webapp2
import datamodel

class MainHandler(webapp2.RequestHandler):
    def get(self):
        new_entry = datamodel.WatchJob(name='test', interval=5, secret='test')
        new_entry.generateSecret()
        new_entry.put()
        self.response.write('Hello world!')

app = webapp2.WSGIApplication([('/', MainHandler)], debug=True)
