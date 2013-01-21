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