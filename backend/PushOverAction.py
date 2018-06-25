import httplib
import logging
import urllib

from google.appengine.ext import db

from datamodels import Action


class PushOverAction(Action):
    token = db.StringProperty(required=True)
    userkeys = db.StringListProperty()
    custom_sound = db.StringProperty(default='')

    def perform_action(self):
        if not self.enabled:
            return

        self.update_performed()

        if self.custom_sound != '':
            sound = self.custom_sound
        else:
            if self.failure_action:
                sound = 'falling'
            else:
                sound = 'bike'

        for user_key in self.userkeys:
            logging.info('PushOverAction sending message')
            conn = httplib.HTTPSConnection('api.pushover.net:443')
            conn.request('POST', '/1/messages.json',
                         urllib.urlencode({
                             'token': self.token,
                             'user': user_key,
                             'message': self.message,
                             'sound': sound
                         }), {'Content-type': 'application/x-www-form-urlencoded'})
            conn.getresponse()
