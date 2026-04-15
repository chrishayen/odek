# Requirement: "a monitoring system with an embedded time-series database"

In-memory time series indexed by metric name and label set, with scrape ingestion, range queries, and a small query function surface.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.http
    std.http.get
      fn (url: string) -> result[string, string]
      + performs an HTTP GET and returns the response body
      - returns error on non-2xx status
      # http
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + returns a 64-bit FNV-1a hash
      # hashing

tsdb
  tsdb.new
    fn () -> tsdb_state
    + creates an empty time series store
    # construction
  tsdb.label_set_id
    fn (labels: map[string, string]) -> u64
    + returns a stable id for a label set, order-independent
    # indexing
    -> std.hash.fnv64
  tsdb.append
    fn (db: tsdb_state, metric: string, labels: map[string, string], value: f64, ts_ms: i64) -> tsdb_state
    + appends a sample to the series identified by (metric, labels)
    + creates the series on first write
    # ingestion
  tsdb.query_range
    fn (db: tsdb_state, metric: string, labels: map[string, string], start_ms: i64, end_ms: i64) -> list[sample]
    + returns samples falling within the inclusive range
    - returns empty list when the series is unknown
    # query
  tsdb.rate
    fn (samples: list[sample], window_ms: i64) -> f64
    + returns per-second rate over the trailing window assuming a counter
    - returns 0 when fewer than two samples fall in the window
    # query
  tsdb.scrape_target
    fn (db: tsdb_state, url: string) -> result[tsdb_state, string]
    + fetches a text-exposition page and appends all parsed samples
    - returns error when HTTP fetch fails
    # ingestion
    -> std.http.get
    -> std.time.now_millis
  tsdb.parse_exposition
    fn (body: string) -> list[tuple[string, map[string, string], f64]]
    + parses the text exposition format into metric samples
    + ignores comment and HELP lines
    # parsing
  tsdb.evaluate_alert
    fn (db: tsdb_state, metric: string, labels: map[string, string], threshold: f64, window_ms: i64) -> bool
    + returns true when the latest rate over the window exceeds the threshold
    # alerting
    -> tsdb.query_range
    -> tsdb.rate
    -> std.time.now_millis
  tsdb.series_list
    fn (db: tsdb_state) -> list[tuple[string, map[string, string]]]
    + returns all known series by (metric, labels)
    # introspection
  tsdb.delete_series
    fn (db: tsdb_state, metric: string, labels: map[string, string]) -> tsdb_state
    + removes the matching series from the store
    # mutation
