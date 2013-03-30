from datetime import tzinfo, timedelta, datetime

class German_tzinfo(tzinfo):
  """Implementation of the German timezone."""

  def utcToLocal(self, dt):
    return dt + self.utcoffset(dt);

  def utcoffset(self, dt):
    return timedelta(hours=+1) + self.dst(dt)

  def _LastSunday(self, dt):
    """First Sunday on or after dt."""
    return dt + timedelta(days=(dt.weekday()-6))

  def dst(self, dt):
    # 2 am on the last Sunday in March
    dst_start = self._LastSunday(datetime(dt.year, 3, 31, 2))

    # 3 am on the last Sunday in October
    dst_end = self._LastSunday(datetime(dt.year, 10, 1, 3))

    if dst_start <= dt.replace(tzinfo=None) < dst_end:
      return timedelta(hours=1)
    else:
      return timedelta(hours=0)