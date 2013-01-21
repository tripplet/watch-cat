#!/usr/bin/env python

import os, hashlib
from google.appengine.ext import db

class PollEvent(db.Model):
  host       = db.StringProperty(required=True)
  port       = db.IntegerProperty(required=True)
  poll_type  = db.StringProperty(default='http', choices=set(['raw', 'http']))

class EmailAction(db.Model):
  enabled    = db.BooleanProperty(required=True, default=False)
  address    = db.EmailProperty(required=True)
  message    = db.StringProperty()

class WebrequestAction(db.Model):
  enabled    = db.BooleanProperty(required=True, default=False)
  url        = db.LinkProperty(required=True)
  message    = db.StringProperty()

class IMAction(db.Model):
  enabled    = db.BooleanProperty(required=True, default=False)
  im         = db.IMProperty(required=True)
  message    = db.StringProperty()


class WatchJob(db.Model):
  name       = db.StringProperty(required=True)
  watch_type = db.StringProperty(required=True, default='push', choices=set(['push', 'poll']))
  enabled    = db.BooleanProperty(required=True, default=False)
  interval   = db.IntegerProperty(required=True)
  created    = db.DateTimeProperty(required=True, auto_now_add=True)
  secret     = db.StringProperty()
  last_fail  = db.DateTimeProperty()
  last_seen  = db.DateTimeProperty()
  last_ip    = db.StringProperty()
  actions    = db.ListProperty(db.Key)
  poll       = db.ReferenceProperty(PollEvent)

  def generateSecret(self):
    self.secret = hashlib.sha1(os.urandom(1024)).hexdigest()