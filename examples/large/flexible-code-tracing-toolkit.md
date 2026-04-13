# Requirement: "a flexible code tracing toolkit"

Instruments function entry, exit, and line events against a flexible predicate language. The std layer exposes low-level tracing hooks; the project layer composes predicates and formats events.

std
  std.tracing
    std.tracing.install_trace_hook
      @ (callback: trace_callback) -> result[hook_id, string]
      + installs a callback invoked on function entry, exit, and line events
      - returns error when another hook is already installed
      # runtime
    std.tracing.remove_trace_hook
      @ (id: hook_id) -> result[void, string]
      + uninstalls the hook
      # runtime
  std.reflection
    std.reflection.frame_function_name
      @ (frame: frame_ref) -> string
      + returns the qualified function name of frame
      # introspection
    std.reflection.frame_filename
      @ (frame: frame_ref) -> string
      + returns the source file of frame
      # introspection
    std.reflection.frame_line
      @ (frame: frame_ref) -> i32
      + returns the current line number in frame
      # introspection
    std.reflection.frame_local
      @ (frame: frame_ref, name: string) -> optional[value]
      + returns a local variable by name
      - returns none when the variable is not in scope
      # introspection

tracer
  tracer.predicate_eq
    @ (field: string, value: string) -> predicate
    + matches events whose field equals value
    # predicates
  tracer.predicate_regex
    @ (field: string, pattern: string) -> result[predicate, string]
    + matches events whose field matches pattern
    - returns error on invalid pattern
    # predicates
  tracer.predicate_and
    @ (left: predicate, right: predicate) -> predicate
    + matches when both predicates match
    # predicates
  tracer.predicate_or
    @ (left: predicate, right: predicate) -> predicate
    + matches when either predicate matches
    # predicates
  tracer.predicate_not
    @ (inner: predicate) -> predicate
    + inverts a predicate
    # predicates
  tracer.describe_event
    @ (kind: event_kind, frame: frame_ref) -> trace_event
    + captures function, filename, line, and kind into a trace_event
    # event_capture
    -> std.reflection.frame_function_name
    -> std.reflection.frame_filename
    -> std.reflection.frame_line
  tracer.event_matches
    @ (event: trace_event, p: predicate) -> bool
    + returns true when the event satisfies the predicate
    # evaluation
  tracer.format_event
    @ (event: trace_event) -> string
    + renders an event as a single log line
    # rendering
  tracer.start
    @ (p: predicate, sink: event_sink) -> result[tracer_state, string]
    + installs a hook that forwards matching events to sink
    - returns error when another tracer is active
    # control
    -> std.tracing.install_trace_hook
  tracer.stop
    @ (state: tracer_state) -> result[void, string]
    + uninstalls the hook and drains remaining events
    # control
    -> std.tracing.remove_trace_hook
  tracer.capture_local
    @ (frame: frame_ref, name: string) -> optional[value]
    + returns a snapshot of a local variable for inclusion in events
    # event_capture
    -> std.reflection.frame_local
