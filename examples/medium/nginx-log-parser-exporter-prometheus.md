# Requirement: "a web server access log parser that exposes aggregated counters and histograms"

Parse combined-format access lines, aggregate by status and path, and expose a text-format metrics snapshot.

std
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits a string by the given separator
      # strings
  std.parse
    std.parse.parse_i64
      @ (s: string) -> result[i64, string]
      + parses a signed integer from text
      - returns error when the string is not a valid integer
      # parsing
    std.parse.parse_f64
      @ (s: string) -> result[f64, string]
      + parses a 64-bit float from text
      - returns error when the string is not a valid float
      # parsing

access_log
  access_log.parse_line
    @ (line: string) -> result[log_entry, string]
    + parses method, path, status, bytes_sent, and request_time
    - returns error when the line has fewer fields than expected
    - returns error when status is not numeric
    # parsing
    -> std.strings.split
    -> std.parse.parse_i64
    -> std.parse.parse_f64
  access_log.new_metrics
    @ () -> metrics_state
    + creates an empty metrics aggregator
    # construction
  access_log.record
    @ (m: metrics_state, entry: log_entry) -> metrics_state
    + increments the status-code counter and updates the request_time histogram
    # aggregation
  access_log.record_line
    @ (m: metrics_state, line: string) -> result[metrics_state, string]
    + parses and records in one step
    - returns error when the line is malformed
    # aggregation
    -> access_log.parse_line
  access_log.counter
    @ (m: metrics_state, status: i32) -> i64
    + returns the count for the given status code
    - returns 0 when the status code was never seen
    # inspection
  access_log.histogram_buckets
    @ (m: metrics_state) -> list[histogram_bucket]
    + returns cumulative buckets for request_time
    # inspection
  access_log.render_text_exposition
    @ (m: metrics_state) -> string
    + emits the metrics in a line-oriented exposition format
    # export
