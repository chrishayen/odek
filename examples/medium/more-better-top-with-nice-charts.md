# Requirement: "an interactive process monitor with charted resource usage"

The library owns sampling, ring-buffer history, and rendering to a text grid. The host supplies the terminal surface.

std
  std.proc
    std.proc.list_processes
      @ () -> list[proc_info]
      + returns a snapshot of running processes with pid, command, cpu_pct, rss_bytes
      # process
  std.system
    std.system.cpu_usage
      @ () -> f64
      + returns aggregate cpu usage as a percentage
      # system
    std.system.memory_usage
      @ () -> tuple[u64, u64]
      + returns (used_bytes, total_bytes)
      # system
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

proc_monitor
  proc_monitor.new
    @ (history_points: i32) -> proc_monitor_state
    + creates a monitor retaining the given number of history samples
    # construction
  proc_monitor.sample
    @ (s: proc_monitor_state) -> proc_monitor_state
    + records current cpu, memory, and per-process stats into the ring buffer
    # sampling
    -> std.system.cpu_usage
    -> std.system.memory_usage
    -> std.proc.list_processes
    -> std.time.now_millis
  proc_monitor.top_processes
    @ (s: proc_monitor_state, n: i32, sort_by: string) -> list[proc_info]
    + returns the top n processes sorted by "cpu" or "memory"
    - returns the empty list when no samples have been taken
    # query
  proc_monitor.history
    @ (s: proc_monitor_state, series: string) -> list[f64]
    + returns the history buffer for "cpu" or "memory"
    # query
  proc_monitor.render_sparkline
    @ (values: list[f64], width: i32) -> string
    + returns a unicode sparkline of the values clipped to width
    - returns an empty string when width <= 0
    # rendering
  proc_monitor.render_dashboard
    @ (s: proc_monitor_state, width: i32, height: i32) -> string
    + returns a text grid containing the top-process table and sparklines
    # rendering
    -> proc_monitor.top_processes
    -> proc_monitor.history
