https://cloud.google.com/appengine/docs/admin-api/accessing-the-api


https://cloud.google.com/appengine/docs/standard/python/config/setting-autoscaling-params-in-explorer?hl=de


curl --request PATCH \
  'https://appengine.googleapis.com/v1/apps/watch-cat/services/default/versions/main?updateMask=automatic_scaling.standard_scheduler_settings.max_instances' \
  --header 'Authorization: Bearer [YOUR_BEARER_TOKEN]' \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --data '{"automaticScaling":{"standardSchedulerSettings":{"maxInstances":1}}}' \
  --compressed

curl --request PATCH \
  'https://appengine.googleapis.com/v1/apps/watch-cat/services/default/versions/main?updateMask=automatic_scaling.standard_scheduler_settings.max_instances&key=[YOUR_API_KEY]' \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --data '{"automaticScaling":{"standardSchedulerSettings":{"maxInstances":1}}}' \
  --compressed