# Requirement: "a resiliency testing tool that injects random instance failures"

Given a population of instances and a schedule, decide which instances to terminate each tick, respecting exclusion windows and per-group rate limits.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.random
    std.random.next_u64
      @ () -> u64
      + returns a uniformly distributed 64-bit value
      # randomness

chaos
  chaos.new_schedule
    @ (groups: list[instance_group]) -> schedule_state
    + creates a schedule over the given instance groups
    # construction
  chaos.is_in_window
    @ (state: schedule_state, group: string) -> bool
    + returns true when the group is currently inside its allowed termination window
    # scheduling
    -> std.time.now_seconds
  chaos.pick_victim
    @ (state: schedule_state, group: string) -> result[string, string]
    + returns a random instance id from the group using a uniform draw
    - returns error when the group is empty
    - returns error when the group's per-tick limit is exhausted
    # selection
    -> std.random.next_u64
  chaos.record_termination
    @ (state: schedule_state, group: string, instance_id: string) -> schedule_state
    + increments the group's tick counter so subsequent picks observe the rate limit
    # bookkeeping
    -> std.time.now_seconds
  chaos.run_tick
    @ (state: schedule_state, terminator: instance_terminator) -> result[list[string], string]
    + walks every in-window group, picks victims up to the limit, and terminates them
    + returns the list of terminated instance ids
    # orchestration
