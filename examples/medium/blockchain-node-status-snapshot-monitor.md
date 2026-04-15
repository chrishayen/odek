# Requirement: "a monitoring library that polls a blockchain node and exposes status snapshots"

The library handles polling, parsing, and snapshot storage. Terminal rendering is the caller's concern.

std
  std.http
    std.http.get
      fn (url: string) -> result[string, string]
      + performs a GET and returns the response body
      - returns error on network failure or non-2xx status
      # http
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

node_monitor
  node_monitor.new
    fn (endpoint: string) -> monitor_state
    + stores the node endpoint and initializes an empty snapshot history
    # construction
  node_monitor.poll
    fn (state: monitor_state) -> result[monitor_state, string]
    + fetches status from the node endpoint and appends a timestamped snapshot
    - returns error when the endpoint is unreachable or returns an unparsable body
    # polling
    -> std.http.get
    -> std.json.parse_object
    -> std.time.now_millis
  node_monitor.latest
    fn (state: monitor_state) -> optional[node_snapshot]
    + returns the most recently recorded snapshot, or none if never polled
    # query
  node_monitor.history
    fn (state: monitor_state) -> list[node_snapshot]
    + returns all retained snapshots in chronological order
    # query
  node_monitor.is_synced
    fn (snapshot: node_snapshot, tip_slot: i64, tolerance_slots: i64) -> bool
    + returns true when the snapshot slot is within tolerance of the network tip
    # health_check
