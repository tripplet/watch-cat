import logging

from google.appengine.api import app_identity
from google.appengine.api import mail
from google.appengine.ext import db
from datamodels import Action


class EmailAction(Action):
    address = db.EmailProperty(required=True)
    subject = db.StringProperty(required=True)

    appengine_mailadress = 'warning@' + app_identity.get_application_id() + '.appspotmail.com'

    def perform_action(self):
        if not self.enabled:
            return

        self.update_performed()
        logging.info('EmailAction sending to:%s', self.address)

        mail.send_mail(sender=EmailAction.appengine_mailadress,
                       to=self.address,
                       subject=self.subject,
                       body=self.message)
