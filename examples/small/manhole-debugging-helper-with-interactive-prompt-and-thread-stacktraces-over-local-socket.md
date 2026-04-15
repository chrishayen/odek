# Requirement: "a debugging helper exposing an interactive prompt and thread stacktraces over a local socket"

The library accepts a connected socket and drives the debug loop; opening and accepting the socket is the caller's job.

std
  std.runtime
    std.runtime.all_thread_stacks
      fn () -> list[thread_stack]
      + returns the current call stack for every live thread
      # runtime

manhole
  manhole.render_stacks
    fn (stacks: list[thread_stack]) -> string
    + formats thread stacks as a multiline string with a header per thread
    + returns "no threads" when stacks is empty
    # formatting
  manhole.handle_command
    fn (session: session_state, line: string) -> tuple[string, session_state]
    + "stacks" returns rendered stacks; "help" returns available commands
    + unknown commands return an error message without exiting
    + "quit" returns a quit marker for the caller to close the session
    # dispatch
    -> std.runtime.all_thread_stacks
    -> manhole.render_stacks
  manhole.new_session
    fn () -> session_state
    + returns a fresh interactive session state
    # construction
