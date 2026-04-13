# Requirement: "a client that routes SQL queries across a cluster of database hosts with primary/replica roles"

Health checks run against injected probe functions so the library stays database-agnostic. Routing picks a host by role; queries are executed through a pluggable execute function.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

dbcluster
  dbcluster.new
    @ (hosts: list[host_config]) -> cluster_state
    + builds a cluster state from a list of host configurations
    ? each host_config carries an address, role hint, and execute function
    # construction
  dbcluster.refresh_health
    @ (state: cluster_state, probe: fn(host_config) -> result[host_role, string]) -> cluster_state
    + probes every host and updates its known role and last-seen timestamp
    + hosts that fail the probe are marked unhealthy
    # health
    -> std.time.now_millis
  dbcluster.pick_primary
    @ (state: cluster_state) -> result[host_config, string]
    + returns a healthy host currently in the primary role
    - returns error when no primary is healthy
    # routing
  dbcluster.pick_replica
    @ (state: cluster_state) -> result[host_config, string]
    + returns a healthy replica, round-robining across available ones
    - returns error when no replica is healthy
    # routing
  dbcluster.execute_write
    @ (state: cluster_state, query: string, params: list[string]) -> result[query_result, string]
    + runs the query on the current primary
    - returns error when no primary is healthy
    - propagates errors from the host's execute function
    # execution
    -> dbcluster.pick_primary
  dbcluster.execute_read
    @ (state: cluster_state, query: string, params: list[string]) -> result[query_result, string]
    + runs the query on a healthy replica, falling back to the primary when none are available
    - returns error when no host is healthy
    # execution
    -> dbcluster.pick_replica
    -> dbcluster.pick_primary
