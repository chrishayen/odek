# Requirement: "a chaos engineering toolkit that orchestrates experiments against running systems"

An experiment is a plan of steady-state probes, faults to inject, and rollback actions. The project layer runs plans; std provides primitives for timing, randomness, and JSON.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      fn (duration: i64) -> void
      + blocks the caller for the given duration
      # time
  std.random
    std.random.uniform
      fn () -> f64
      + returns a uniformly distributed value in [0, 1)
      # random
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a flat JSON object
      - returns error on malformed input
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a map as JSON
      # serialization

chaos_toolkit
  chaos_toolkit.new_experiment
    fn (name: string) -> experiment_state
    + creates an empty experiment with the given name
    # construction
  chaos_toolkit.add_steady_state_probe
    fn (state: experiment_state, probe: probe_spec) -> experiment_state
    + attaches a probe that will be verified before and after the experiment
    # configuration
  chaos_toolkit.add_action
    fn (state: experiment_state, action: action_spec) -> experiment_state
    + appends a fault-inducing action to the method phase
    # configuration
  chaos_toolkit.add_rollback
    fn (state: experiment_state, action: action_spec) -> experiment_state
    + appends a compensating action to run during the rollback phase
    # configuration
  chaos_toolkit.run
    fn (state: experiment_state, executor: fn(action_spec) -> result[void, string]) -> experiment_report
    + runs steady-state checks, then the method phase, then another steady-state check, then rollbacks
    + returns a report containing each step's outcome and timing
    - aborts the method phase on the first action error and runs rollbacks
    # orchestration
    -> std.time.now_millis
  chaos_toolkit.verify_steady_state
    fn (state: experiment_state, executor: fn(probe_spec) -> result[map[string, string], string]) -> result[void, string]
    + runs every probe and returns ok only when all match their tolerances
    - returns error naming the first probe whose result falls outside tolerance
    # verification
  chaos_toolkit.sample_with_probability
    fn (p: f64) -> bool
    + returns true with probability p; used to flip faults on or off
    # sampling
    -> std.random.uniform
  chaos_toolkit.encode_plan
    fn (state: experiment_state) -> string
    + serializes an experiment to a canonical JSON string
    # serialization
    -> std.json.encode_object
  chaos_toolkit.decode_plan
    fn (raw: string) -> result[experiment_state, string]
    + restores an experiment from a JSON string
    - returns error on malformed JSON
    # serialization
    -> std.json.parse_object
  chaos_toolkit.pace_steps
    fn (state: experiment_state, delay_millis: i64) -> experiment_state
    + configures a sleep between consecutive actions in the method phase
    # configuration
    -> std.time.sleep_millis
