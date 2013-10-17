from datetime import tzinfo, timedelta, datetime

class German_tzinfo(tzinfo):
  """Implementation of the German timezone."""

  def utcToLocal(self, dt):
    return dt + self.utcoffset(dt)

  def utcoffset(self, dt):
    return timedelta(hours=+1) + self.dst(dt)

  def _LastSunday(self, dt):
    """Previous Sunday on or bevor dt."""
    if dt.isoweekday() == 7:
      return dt
    else:
      return dt + timedelta(days=-dt.isoweekday())

  def dst(self, dt):
    # 2 am (02:00) on the last Sunday in March
    dst_start = self._LastSunday(datetime(dt.year, 3, 31, 2))

    # 3 am (03:00) on the last Sunday in October
    dst_end = self._LastSunday(datetime(dt.year, 10, 31, 3))

    if dst_start <= dt.replace(tzinfo=None) < dst_end:
      return timedelta(hours=+1)
    else:
      return timedelta(hours=0)