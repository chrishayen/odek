# Requirement: "a supervision tree library with configurable restart strategies"

Models a supervisor as a tree of child specs. When a child fails, the library applies a restart strategy (one-for-one, one-for-all, rest-for-one) and reports which children must restart.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

supervisor
  supervisor.new_child_spec
    @ (id: string, start: function[void, child_state], restart: string) -> child_spec
    + builds a child spec with an id, start function, and restart policy string ("permanent", "transient", or "temporary")
    # configuration
  supervisor.new
    @ (strategy: string, children: list[child_spec]) -> result[sup_state, string]
    + builds a supervisor with a named strategy and an ordered list of children
    - returns error when strategy is not "one_for_one", "one_for_all", or "rest_for_one"
    - returns error when child ids are not unique
    # configuration
  supervisor.start_all
    @ (sup: sup_state) -> sup_state
    + starts every child in order and stores their running state
    # lifecycle
  supervisor.handle_exit
    @ (sup: sup_state, failed_id: string, reason: string) -> result[tuple[sup_state, list[string]], string]
    + applies the supervisor's strategy and returns the new state plus the list of child ids that were restarted
    - returns error when failed_id is not a child
    # lifecycle
    -> std.time.now_millis
  supervisor.intensity_exceeded
    @ (sup: sup_state, max_restarts: i32, window_ms: i64) -> bool
    + returns true when restart events within the window exceed max_restarts
    # health
    -> std.time.now_millis
  supervisor.stop_all
    @ (sup: sup_state) -> sup_state
    + stops every child in reverse start order
    # lifecycle
