#!/bin/sh

API_KEY=test
APP=watch-cat

gcloud app deploy --version main
curl --request PATCH \
  'https://appengine.googleapis.com/v1/apps/${APP}/services/default/versions/main?updateMask=automaticScaling.max_idle_instances&key=${API_KEY}' \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --data '{"automaticScaling":{"standardSchedulerSettings":{"maxInstances":1}}}' \
  --compressed
