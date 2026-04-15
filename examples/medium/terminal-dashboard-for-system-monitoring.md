# Requirement: "a terminal dashboard for system monitoring"

Collects system metrics (cpu, memory, disk, network) and renders them as a text dashboard. Metric collection and terminal size go through std primitives; the project layer composes samples into a displayable frame.

std
  std.system
    std.system.cpu_usage
      fn () -> f64
      + returns overall cpu utilization in [0, 1]
      # system
    std.system.memory_usage
      fn () -> tuple[i64, i64]
      + returns (used_bytes, total_bytes)
      # system
    std.system.disk_usage
      fn (path: string) -> result[tuple[i64, i64], string]
      + returns (used_bytes, total_bytes) for the mount containing the path
      - returns error when the path is not mounted
      # system
    std.system.net_bytes
      fn () -> tuple[i64, i64]
      + returns cumulative (rx_bytes, tx_bytes) across all interfaces
      # system
  std.term
    std.term.size
      fn () -> tuple[i32, i32]
      + returns (columns, rows) of the controlling terminal
      # terminal
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

sysdash
  sysdash.sample
    fn () -> metric_snapshot
    + captures a single snapshot of cpu, memory, disk, and network counters
    # collection
    -> std.system.cpu_usage
    -> std.system.memory_usage
    -> std.system.disk_usage
    -> std.system.net_bytes
    -> std.time.now_millis
  sysdash.net_rate
    fn (previous: metric_snapshot, current: metric_snapshot) -> tuple[f64, f64]
    + returns (rx_bytes_per_sec, tx_bytes_per_sec) between two snapshots
    ? returns (0, 0) when the elapsed time is zero
    # derivation
  sysdash.render_frame
    fn (snapshot: metric_snapshot, rx_rate: f64, tx_rate: f64, width: i32, height: i32) -> string
    + produces a text frame fit to the given size with bars and values
    # rendering
  sysdash.format_bytes
    fn (value: i64) -> string
    + returns a human-friendly size like "1.2 GiB"
    # formatting
  sysdash.format_percent
    fn (ratio: f64) -> string
    + returns a percent string like "42.0%"
    # formatting
  sysdash.default_size
    fn () -> tuple[i32, i32]
    + returns the current terminal dimensions
    # layout
    -> std.term.size
