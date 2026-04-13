# Requirement: "a library for managing application shutdown with signals"

Collects shutdown hooks and runs them when a termination signal arrives. Signal delivery goes through a thin std primitive so tests can drive shutdown deterministically.

std
  std.signal
    std.signal.wait_for_termination
      @ () -> string
      + blocks until a termination signal is delivered and returns its name
      # signals

shutdown
  shutdown.new
    @ () -> shutdown_state
    + creates an empty shutdown manager with no registered hooks
    # construction
  shutdown.on_close
    @ (state: shutdown_state, name: string, hook: fn() -> result[void, string]) -> shutdown_state
    + registers a named shutdown hook to run when termination is triggered
    # registration
  shutdown.wait
    @ (state: shutdown_state) -> result[void, string]
    + blocks until a termination signal, then runs each hook in registration order
    - returns an aggregated error describing every hook that failed
    # lifecycle
    -> std.signal.wait_for_termination
  shutdown.trigger
    @ (state: shutdown_state, reason: string) -> result[void, string]
    + runs hooks immediately without waiting for a signal
    - returns an aggregated error describing every hook that failed
    # lifecycle
