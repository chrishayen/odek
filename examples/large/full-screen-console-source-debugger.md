# Requirement: "a full-screen console-based source-level debugger"

A debugger library: stepping, breakpoints, variable inspection, and a text-UI renderer. The caller drives the event loop and supplies a language runtime adapter.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the full contents of a text file
      - returns error when the file does not exist
      # filesystem
  std.term
    std.term.clear_screen
      fn () -> string
      + returns the ANSI sequence that clears the terminal
      # terminal
    std.term.move_cursor
      fn (row: i32, col: i32) -> string
      + returns the ANSI sequence that moves the cursor
      # terminal
    std.term.read_key
      fn () -> result[string, string]
      + returns the next key press as a symbolic name ("up", "q", "space")
      - returns error when stdin is closed
      # terminal
    std.term.size
      fn () -> tuple[i32, i32]
      + returns (rows, cols) of the current terminal
      # terminal

debugger
  debugger.new_session
    fn (source_path: string) -> result[session_state, string]
    + creates a debug session and loads source text for the given file
    - returns error when the source file cannot be read
    # construction
    -> std.fs.read_all
  debugger.set_breakpoint
    fn (session: session_state, path: string, line: i32) -> session_state
    + records a breakpoint at the given file and line
    # breakpoints
  debugger.clear_breakpoint
    fn (session: session_state, path: string, line: i32) -> session_state
    + removes a breakpoint matching file and line
    - leaves the session unchanged when no matching breakpoint exists
    # breakpoints
  debugger.list_breakpoints
    fn (session: session_state) -> list[breakpoint_record]
    + returns all currently set breakpoints
    # breakpoints
  debugger.step_into
    fn (session: session_state, runtime: runtime_adapter) -> result[session_state, string]
    + advances the program counter into the next callable
    - returns error when the runtime reports the program has terminated
    # stepping
  debugger.step_over
    fn (session: session_state, runtime: runtime_adapter) -> result[session_state, string]
    + advances past the current statement without descending
    # stepping
  debugger.continue_exec
    fn (session: session_state, runtime: runtime_adapter) -> result[session_state, string]
    + runs until the next breakpoint or program termination
    # stepping
  debugger.current_frame
    fn (session: session_state) -> optional[stack_frame]
    + returns the topmost stack frame when the program is paused
    - returns none when the program is running or ended
    # inspection
  debugger.local_variables
    fn (session: session_state, frame: stack_frame) -> map[string, string]
    + returns local variable name to display-string map for a frame
    # inspection
  debugger.evaluate
    fn (session: session_state, frame: stack_frame, expression: string) -> result[string, string]
    + evaluates an expression in the context of the given frame
    - returns error when the expression is invalid in the host runtime
    # inspection
  debugger.render_view
    fn (session: session_state) -> string
    + returns a full-screen text render with source, stack, and locals panes
    # ui
    -> std.term.clear_screen
    -> std.term.move_cursor
    -> std.term.size
  debugger.handle_key
    fn (session: session_state, key: string) -> session_state
    + maps a key name to a session action and returns the updated session
    + unknown keys leave the session unchanged
    # ui
    -> std.term.read_key
