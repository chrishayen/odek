# Requirement: "a library that collects system metrics and optionally relays them to a pluggable time-series sink"

Samples CPU, memory, and disk usage at intervals and ships them to a sink adapter.

std
  std.os
    std.os.cpu_usage_percent
      fn () -> f64
      + returns aggregate CPU utilization since last call
      # system
    std.os.mem_used_bytes
      fn () -> i64
      + returns resident memory in use
      # system
    std.os.disk_used_bytes
      fn (mount: string) -> result[i64, string]
      + returns used bytes on the given mount
      - returns error when mount is not present
      # system
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

sysmetrics
  sysmetrics.sample
    fn () -> map[string, f64]
    + returns a point-in-time map of metric name to value
    # sampling
    -> std.os.cpu_usage_percent
    -> std.os.mem_used_bytes
  sysmetrics.format_text
    fn (sample: map[string, f64], timestamp_ms: i64) -> string
    + renders a sample as newline-delimited "name value timestamp" lines
    # formatting
  sysmetrics.format_json
    fn (sample: map[string, f64], timestamp_ms: i64) -> string
    + renders a sample as a JSON object with a timestamp field
    # formatting
  sysmetrics.relay
    fn (sample: map[string, f64], sink: bus_state) -> result[void, string]
    + forwards a sample to the attached sink
    - returns error when the sink rejects the write
    ? sink is a pluggable adapter opened by the caller
    # sink
    -> std.time.now_millis
  sysmetrics.disk_sample
    fn (mounts: list[string]) -> result[map[string, f64], string]
    + returns per-mount used bytes
    - returns error if any mount cannot be read
    # sampling
    -> std.os.disk_used_bytes
