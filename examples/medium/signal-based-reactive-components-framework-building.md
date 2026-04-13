# Requirement: "a signal-based reactive component framework for text user interfaces"

Signals track dependencies, components re-render when the signals they read change.

std
  std.term
    std.term.clear_screen
      @ () -> void
      + writes the clear-screen escape sequence to stdout
      # terminal
    std.term.move_cursor
      @ (row: i32, col: i32) -> void
      + positions the cursor at the given row and column
      # terminal

reactive_tui
  reactive_tui.signal_new
    @ (initial: string) -> signal_handle
    + creates a signal holding the given value
    # reactivity
  reactive_tui.signal_get
    @ (sig: signal_handle) -> string
    + returns the current value and records a dependency on the active tracker
    # reactivity
  reactive_tui.signal_set
    @ (sig: signal_handle, value: string) -> void
    + updates the value and schedules all dependent subscribers for re-evaluation
    # reactivity
  reactive_tui.effect
    @ (fn: effect_fn) -> effect_handle
    + runs fn once, tracks which signals it reads, and re-runs it when any of them change
    # effects
  reactive_tui.component
    @ (render: render_fn) -> component_handle
    + wraps a render function as a component tied to its captured signals
    # components
  reactive_tui.mount
    @ (root: component_handle) -> void
    + performs an initial render and enters the event loop until the root is disposed
    # lifecycle
    -> std.term.clear_screen
    -> std.term.move_cursor
  reactive_tui.dispose
    @ (comp: component_handle) -> void
    + unsubscribes the component from all its signals and removes it from the tree
    # lifecycle
