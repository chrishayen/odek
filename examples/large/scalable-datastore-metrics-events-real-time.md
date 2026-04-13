# Requirement: "a time-series datastore for metrics and events with range queries and simple aggregations"

An append-only per-series log with indexed timestamps, supporting ingest, range queries, and a small set of aggregations.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.fs
    std.fs.append
      @ (path: string, data: bytes) -> result[void, string]
      + appends data to a file, creating it if missing
      - returns error on write failure
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file's full content
      - returns error when the file cannot be read
      # filesystem
  std.encoding
    std.encoding.varint_encode
      @ (n: i64) -> bytes
      + encodes an integer as unsigned varint
      # encoding
    std.encoding.varint_decode
      @ (data: bytes, offset: i32) -> result[tuple[i64,i32], string]
      + decodes a varint and returns the value plus next offset
      - returns error on truncated input
      # encoding

tsstore
  tsstore.point
    @ (timestamp_ms: i64, value: f64, tags: map[string,string]) -> point
    + constructs a single data point
    # construction
  tsstore.series_id
    @ (metric: string, tags: map[string,string]) -> string
    + returns a canonical series identifier by sorting tag keys
    # identity
  tsstore.open
    @ (data_dir: string) -> result[tsstore_state, string]
    + opens or creates a store rooted at data_dir
    # construction
  tsstore.ingest
    @ (state: tsstore_state, metric: string, p: point) -> result[tsstore_state, string]
    + appends a point to the series log, updating in-memory index
    - returns error when the point's timestamp is older than the series' latest
    # ingest
    -> std.fs.append
    -> std.encoding.varint_encode
  tsstore.range_query
    @ (state: tsstore_state, metric: string, tags: map[string,string], from_ms: i64, to_ms: i64) -> result[list[point], string]
    + returns all points in the window, sorted by timestamp
    - returns error when the series does not exist
    # query
    -> std.fs.read_all
    -> std.encoding.varint_decode
  tsstore.aggregate_sum
    @ (points: list[point]) -> f64
    + returns the sum of values
    # aggregation
  tsstore.aggregate_avg
    @ (points: list[point]) -> f64
    + returns the arithmetic mean
    + returns 0.0 on an empty list
    # aggregation
  tsstore.aggregate_min_max
    @ (points: list[point]) -> tuple[f64,f64]
    + returns (min, max) across the points
    ? on empty input returns (0.0, 0.0)
    # aggregation
  tsstore.downsample
    @ (points: list[point], bucket_ms: i64) -> list[point]
    + groups points into fixed-width time buckets and emits the average per bucket
    # aggregation
  tsstore.publish_event
    @ (state: tsstore_state, category: string, message: string) -> result[tsstore_state, string]
    + appends an event record tagged with the current time
    # events
    -> std.time.now_millis
    -> std.fs.append
  tsstore.events_since
    @ (state: tsstore_state, category: string, since_ms: i64) -> result[list[tuple[i64,string]], string]
    + returns events at or after the given time
    - returns error when the category has no log
    # events
    -> std.fs.read_all
