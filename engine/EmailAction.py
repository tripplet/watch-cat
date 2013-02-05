from google.appengine.api import mail
from google.appengine.ext import db

from datamodels import Action

appengine_mailadress = ''

class EmailAction(Action):
  address = db.EmailProperty(required=True)
  subject = db.StringProperty(required=True)

  def performAction(self):
    if not self.enabled:
      return

    self.updatePerformed()
    mail.send_mail(sender  = appengine_mailadress,
                   to      = self.address,
                   subject = self.subject,
                   body    = self.message)