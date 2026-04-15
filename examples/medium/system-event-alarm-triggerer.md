# Requirement: "a library to fire alarms based on system events"

Registers alarm rules on metrics like CPU load, memory, and disk usage, then evaluates samples against those rules and emits alarms when thresholds are crossed.

std
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns the current unix time in seconds
      # time
  std.sysinfo
    std.sysinfo.cpu_percent
      fn () -> f64
      + returns current total CPU usage as a percentage
      # system
    std.sysinfo.memory_used_bytes
      fn () -> i64
      + returns current used memory in bytes
      # system
    std.sysinfo.disk_used_percent
      fn (path: string) -> result[f64, string]
      + returns the percentage of space used at path
      - returns error when path is not a mount point
      # system

alarms
  alarms.new
    fn () -> alarm_registry
    + creates an empty alarm registry
    # construction
  alarms.add_rule
    fn (registry: alarm_registry, rule: alarm_rule) -> alarm_registry
    + returns a registry with rule appended
    ? rule describes metric name, comparison, threshold, and hysteresis
    # configuration
  alarms.sample_metrics
    fn () -> metric_snapshot
    + captures a snapshot of CPU, memory, and disk metrics
    # sampling
    -> std.sysinfo.cpu_percent
    -> std.sysinfo.memory_used_bytes
    -> std.sysinfo.disk_used_percent
    -> std.time.now_seconds
  alarms.evaluate_rule
    fn (rule: alarm_rule, snapshot: metric_snapshot, last_state: alarm_state) -> alarm_state
    + returns alarm_state tracking triggered, recovered, or unchanged
    + respects hysteresis so flapping values do not retrigger
    # evaluation
  alarms.evaluate_all
    fn (registry: alarm_registry, snapshot: metric_snapshot, states: map[string, alarm_state]) -> list[alarm_event]
    + returns a list of fire and recover events produced this tick
    # evaluation
  alarms.format_event
    fn (event: alarm_event) -> string
    + renders an alarm event as a human-readable line
    # rendering
