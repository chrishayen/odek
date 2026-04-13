# Requirement: "a library for registering and running application shutdown hooks in response to termination signals"

Hooks are registered with a priority and run in order when a terminating signal is received.

std
  std.signal
    std.signal.wait_any
      @ (signals: list[string]) -> string
      + blocks until one of the named signals is delivered and returns its name
      # signals

shutdown
  shutdown.new
    @ () -> shutdown_state
    + creates an empty hook registry
    # construction
  shutdown.register
    @ (state: shutdown_state, priority: i32, name: string, hook: hook_fn) -> shutdown_state
    + adds a hook to the registry; lower priority numbers run first
    # registration
  shutdown.unregister
    @ (state: shutdown_state, name: string) -> shutdown_state
    + removes a hook by name; does nothing if not present
    # registration
  shutdown.run_all
    @ (state: shutdown_state) -> list[hook_result]
    + runs every registered hook in priority order, collecting results
    + continues running later hooks even when an earlier one fails
    # execution
  shutdown.wait_and_run
    @ (state: shutdown_state, signals: list[string]) -> list[hook_result]
    + blocks until any listed signal arrives, then runs all hooks
    # orchestration
    -> std.signal.wait_any
