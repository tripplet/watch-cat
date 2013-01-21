#!/usr/bin/env python

import webapp2
import datamodel

class JobHandler(webapp2.RequestHandler):
    def get(self, key):
        self.response.write('Key: ')
        self.response.write(key)


# Main handler ######################################################################
# ###################################################################################
app = webapp2.WSGIApplication([('/job/(\w*)', JobHandler)], debug=True)
# ###################################################################################
# ###################################################################################