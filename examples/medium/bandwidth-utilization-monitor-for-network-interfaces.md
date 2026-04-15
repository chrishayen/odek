# Requirement: "a bandwidth utilization monitor for network interfaces"

Reads interface byte counters at intervals and computes per-process and per-connection throughput.

std
  std.net
    std.net.list_interfaces
      fn () -> result[list[interface_info], string]
      + returns name, index, and byte counters for each interface
      - returns error when the platform has no interface API
      # networking
    std.net.list_connections
      fn () -> result[list[connection_info], string]
      + returns active TCP/UDP connections with pid, local and remote addresses
      # networking
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

bandwidth_monitor
  bandwidth_monitor.snapshot
    fn () -> result[monitor_snapshot, string]
    + captures a point-in-time view of all interfaces and connections
    - returns error when the OS denies permission
    # sampling
    -> std.net.list_interfaces
    -> std.net.list_connections
    -> std.time.now_millis
  bandwidth_monitor.throughput_between
    fn (prev: monitor_snapshot, curr: monitor_snapshot) -> throughput_report
    + returns bytes/sec per interface and per connection over the interval
    - returns zero rates when the timestamps are equal
    # throughput
  bandwidth_monitor.top_talkers
    fn (report: throughput_report, n: i32) -> list[connection_throughput]
    + returns the n connections with the highest combined rx+tx rate
    # analysis
  bandwidth_monitor.group_by_process
    fn (report: throughput_report) -> map[i32, process_throughput]
    + aggregates per-connection rates by pid
    # analysis
