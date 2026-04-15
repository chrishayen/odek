# Requirement: "a library for lightweight programmable telemetry for servers and workstations"

Collects host metrics on a schedule and forwards samples to a pluggable sink. Queries are user-defined expressions.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.host
    std.host.cpu_percent
      fn () -> f64
      + returns current cpu utilization as a percentage in [0, 100]
      # host_metrics
    std.host.memory_used_bytes
      fn () -> i64
      + returns resident memory used by the host in bytes
      # host_metrics
    std.host.hostname
      fn () -> string
      + returns the host identifier for labeling samples
      # host_metrics

telemetry
  telemetry.new_agent
    fn (interval_ms: i64) -> agent_state
    + creates an agent with an empty query set and the given collection interval
    # construction
  telemetry.register_query
    fn (state: agent_state, name: string, expr: string) -> result[agent_state, string]
    + adds a named query whose expression references host metric identifiers
    - returns error when name is empty or already registered
    # query_registration
  telemetry.collect_once
    fn (state: agent_state) -> list[sample]
    + evaluates every registered query and returns one sample per query
    + each sample carries name, numeric value, hostname, and collection timestamp
    # collection
    -> std.time.now_millis
    -> std.host.cpu_percent
    -> std.host.memory_used_bytes
    -> std.host.hostname
  telemetry.eval_expr
    fn (expr: string, metrics: map[string, f64]) -> result[f64, string]
    + resolves bare identifiers against the metrics map and supports +, -, *, /
    - returns error when an identifier is not in metrics
    - returns error on division by zero
    # expression_eval
  telemetry.forward
    fn (samples: list[sample], sink: sample_sink) -> result[void, string]
    + pushes every sample through the sink in order
    - returns error on the first sink failure and stops
    # forwarding
