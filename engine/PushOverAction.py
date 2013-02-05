import httplib, urllib
import logging

from datamodels import Action
from google.appengine.ext import db


class PushOverAction(Action):
  token    = db.StringProperty(required=True)
  userkeys = db.StringListProperty()

  def performAction(self):
    if not self.enabled:
      return

    self.updatePerformed()

    for user_key in self.userkeys:
      logging.info('PushOverAction::performAction to:%s', user_key)
      conn = httplib.HTTPSConnection("api.pushover.net:443")
      conn.request("POST", "/1/messages.json",
        urllib.urlencode({
          "token": self.token,
          "user": user_key,
          "message": self.message,
        }), { "Content-type": "application/x-www-form-urlencoded" })
      conn.getresponse()