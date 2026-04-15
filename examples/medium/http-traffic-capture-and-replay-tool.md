# Requirement: "a tool for capturing http traffic and replaying it against another environment"

Records live requests to a buffer and replays them against a target url, optionally time-scaled.

std
  std.http
    std.http.send
      fn (method: string, url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
      + sends the request and returns the response
      - returns error on connection failure
      # http
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      fn (millis: i64) -> void
      + blocks for the given number of milliseconds
      # time

traffic_replay
  traffic_replay.new_buffer
    fn () -> capture_buffer
    + returns an empty buffer
    # construction
  traffic_replay.capture
    fn (buf: capture_buffer, req: http_request) -> capture_buffer
    + appends a captured request tagged with its arrival timestamp
    # capture
    -> std.time.now_millis
  traffic_replay.replay
    fn (buf: capture_buffer, target_url: string, speedup: f64) -> replay_report
    + sends every captured request to the target, preserving inter-arrival gaps divided by speedup
    + returns a report with one entry per request containing status and latency
    - records an error entry when a send fails
    # replay
    -> std.http.send
    -> std.time.sleep_millis
  traffic_replay.filter
    fn (buf: capture_buffer, path_prefix: string) -> capture_buffer
    + returns a buffer containing only captures whose path starts with the prefix
    # filtering
