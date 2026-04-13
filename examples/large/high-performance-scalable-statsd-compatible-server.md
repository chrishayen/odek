# Requirement: "a scalable statsd-compatible metrics server"

Accepts statsd text packets, aggregates metrics in memory, and flushes to a pluggable sink on an interval.

std
  std.net
    std.net.udp_listen
      @ (addr: string, port: u16) -> result[udp_socket, string]
      + binds a UDP socket to the given address and port
      - returns error when the address is already in use
      # networking
    std.net.udp_recv
      @ (sock: udp_socket) -> result[bytes, string]
      + returns the next datagram payload
      - returns error when the socket is closed
      # networking
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.sync
    std.sync.mutex_new
      @ () -> mutex_handle
      + creates an unlocked mutex
      # concurrency
    std.sync.mutex_with
      @ (m: mutex_handle, body: fn() -> void) -> void
      + runs body while holding the mutex
      # concurrency

statsd_server
  statsd_server.parse_line
    @ (line: string) -> result[metric, string]
    + parses "name:value|type" and optional "|@sample_rate"
    + recognizes counter (c), gauge (g), timer (ms), and set (s) types
    - returns error on missing type suffix or non-numeric value
    # parsing
  statsd_server.parse_packet
    @ (payload: bytes) -> list[metric]
    + splits payload on newlines and parses each line, skipping invalid lines
    # parsing
    -> statsd_server.parse_line
  statsd_server.new_aggregator
    @ () -> aggregator_state
    + creates an empty aggregator holding counters, gauges, timers, and sets
    # construction
  statsd_server.ingest
    @ (agg: aggregator_state, m: metric) -> void
    + counters are summed, gauges are replaced, timers are appended, sets collect unique members
    + applies sample_rate scaling to counters and timers
    # aggregation
  statsd_server.snapshot_and_reset
    @ (agg: aggregator_state) -> snapshot
    + returns the current aggregate values and resets counters, timers, and sets
    + gauges are retained across snapshots
    # aggregation
  statsd_server.summarize_timers
    @ (samples: list[f64]) -> timer_summary
    + returns count, min, max, mean, and p50/p90/p95/p99 percentiles
    - returns zeroed summary when samples is empty
    # statistics
  statsd_server.render_snapshot
    @ (snap: snapshot, prefix: string, now_ms: i64) -> list[metric_point]
    + flattens counters, gauges, timer summaries, and set cardinalities into named points
    # rendering
    -> statsd_server.summarize_timers
  statsd_server.serve
    @ (addr: string, port: u16, sink: fn(list[metric_point]) -> void, flush_ms: i64) -> result[void, string]
    + listens on UDP, aggregates metrics, and calls sink every flush_ms with the snapshot
    - returns error when the socket cannot be bound
    # server
    -> std.net.udp_listen
    -> std.net.udp_recv
    -> std.time.now_millis
    -> std.sync.mutex_new
    -> std.sync.mutex_with
    -> statsd_server.parse_packet
    -> statsd_server.new_aggregator
    -> statsd_server.ingest
    -> statsd_server.snapshot_and_reset
    -> statsd_server.render_snapshot
