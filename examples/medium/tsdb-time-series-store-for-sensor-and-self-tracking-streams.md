# Requirement: "a time-series store for sensor and self-tracking streams"

Create streams, append timestamped readings, query by range, aggregate over windows.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

tsdb
  tsdb.create_stream
    fn (name: string, unit: string) -> result[stream_id, string]
    + creates a named stream with a unit label
    - returns error when a stream with the same name already exists
    # stream_management
  tsdb.append
    fn (stream: stream_id, value: f64) -> result[void, string]
    + appends a reading stamped with the current time
    # write_path
    -> std.time.now_millis
  tsdb.append_at
    fn (stream: stream_id, timestamp_ms: i64, value: f64) -> result[void, string]
    + appends a reading with an explicit timestamp
    - returns error when timestamp is older than the last appended reading
    # write_path
  tsdb.query_range
    fn (stream: stream_id, start_ms: i64, end_ms: i64) -> list[reading]
    + returns readings whose timestamp falls in [start, end)
    + returns empty when no readings match
    # read_path
  tsdb.latest
    fn (stream: stream_id) -> optional[reading]
    + returns the most recent reading
    + returns none when the stream is empty
    # read_path
  tsdb.aggregate_mean
    fn (readings: list[reading]) -> f64
    + returns the arithmetic mean of the values
    - returns 0 when the list is empty
    # aggregation
  tsdb.aggregate_buckets
    fn (readings: list[reading], bucket_ms: i64) -> list[bucket]
    + returns fixed-width buckets with mean value per window
    - returns empty when bucket_ms is not positive
    # aggregation
