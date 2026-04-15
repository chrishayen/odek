# Requirement: "a library that waits for multiple services to become available"

Polls a set of endpoints until each reports healthy or a deadline is reached.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      fn (millis: i64) -> void
      + blocks for the given number of milliseconds
      # time
  std.net
    std.net.tcp_probe
      fn (host: string, port: i32, timeout_ms: i32) -> bool
      + returns true when a tcp connection succeeds within the timeout
      - returns false on refusal, reset, or timeout
      # networking

wait_for
  wait_for.check_all
    fn (targets: list[probe_target]) -> list[bool]
    + returns one boolean per target indicating current reachability
    # probing
    -> std.net.tcp_probe
  wait_for.wait_until_ready
    fn (targets: list[probe_target], timeout_ms: i64, interval_ms: i64) -> result[void, list[string]]
    + returns ok when every target becomes reachable before the deadline
    - returns error with the hostnames still unreachable at timeout
    # polling
    -> std.time.now_millis
    -> std.time.sleep_millis
