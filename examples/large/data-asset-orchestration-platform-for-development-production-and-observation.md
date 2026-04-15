# Requirement: "a data asset orchestration platform for development, production, and observation"

Users declare data assets and the dependencies between them. The library resolves a run order, executes assets, captures materialization events, and exposes a view of the resulting graph and recent runs.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.log
    std.log.info
      fn (message: string) -> void
      + writes an info-level message to the configured sink
      # logging
    std.log.error
      fn (message: string) -> void
      + writes an error-level message to the configured sink
      # logging
  std.graph
    std.graph.topological_sort
      fn (nodes: list[string], edges: list[tuple[string, string]]) -> result[list[string], string]
      + returns nodes in dependency order
      - returns error on cycle, naming the offending edge
      # graph_algorithms

assets
  assets.new_registry
    fn () -> registry_state
    + returns an empty asset registry
    # construction
  assets.define_asset
    fn (reg: registry_state, name: string, deps: list[string], compute: function[map[string, bytes], bytes]) -> result[registry_state, string]
    + registers an asset with its upstream dependency names and compute function
    - returns error when the name is already defined
    - returns error when a declared dependency refers to an undefined asset
    # registration
  assets.plan_run
    fn (reg: registry_state, targets: list[string]) -> result[list[string], string]
    + returns the set of assets that must materialize to satisfy targets, in dependency order
    - returns error when a target is not registered
    - returns error when the dependency graph contains a cycle
    # planning
    -> std.graph.topological_sort
  assets.materialize
    fn (reg: registry_state, plan: list[string]) -> result[run_record, string]
    + executes the plan in order, passing upstream outputs into downstream compute functions
    + records start and end times for each asset
    - returns error when an asset's compute function raises
    # execution
    -> std.time.now_millis
    -> std.log.info
    -> std.log.error
  assets.run_record_for
    fn (run: run_record, asset: string) -> optional[asset_event]
    + returns the materialization event for a named asset in a run
    - returns empty when the asset did not run
    # observability
  assets.list_assets
    fn (reg: registry_state) -> list[string]
    + returns all registered asset names
    # introspection
  assets.dependency_graph
    fn (reg: registry_state) -> list[tuple[string, string]]
    + returns the list of (downstream, upstream) edges for visualization
    # introspection
  assets.recent_runs
    fn (reg: registry_state, limit: i32) -> list[run_record]
    + returns up to limit most recent runs, newest first
    # observability
  assets.invalidate
    fn (reg: registry_state, name: string) -> result[registry_state, string]
    + marks an asset and its downstream dependents as needing re-materialization
    - returns error when the name is not registered
    # lifecycle
