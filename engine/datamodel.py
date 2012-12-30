#!/usr/bin/env python

import datetime
from google.appengine.ext import db

class PollEvent(db.Model):
  host       = db.StringProperty(required=True)
  port       = db.IntegerProperty(required=True)
  poll_type  = db.StringProperty(default='http', choices=set(['raw', 'http']))

class EmailAction(db.Model):
  address    = db.EmailProperty(required=True)
  enabled    = db.BooleanProperty(required=True, default=False)
  message    = db.StringProperty()

class WebrequestAction(db.Model):
  url        = db.LinkProperty(required=True)
  enabled    = db.BooleanProperty(required=True, default=False)
  message    = db.StringProperty()

class IMAction(db.Model):
  im         = db.IMProperty(required=True)
  enabled    = db.BooleanProperty(required=True, default=False)
  message    = db.StringProperty()


class WatchJob(db.Model):
  name       = db.StringProperty(required=True)
  watch_type = db.StringProperty(required=True, default='push', choices=set(['push', 'poll']))
  enabled    = db.BooleanProperty(required=True, default=False)
  interval   = db.IntegerProperty(required=True)
  created    = db.DateTimeProperty(required=True, auto_now_add=True)
  secret     = db.StringProperty()
  last_fail  = db.DateTimeProperty()
  lastseen   = db.DateTimeProperty()
  actions    = db.ListProperty(db.Key)
  poll       = db.ReferenceProperty(PollEvent)