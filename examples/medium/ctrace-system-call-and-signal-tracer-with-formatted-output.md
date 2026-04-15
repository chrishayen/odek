# Requirement: "a library for tracing system calls and signals with formatted output"

Attaches to a process, captures syscall entry/exit events, and formats them as human-readable records.

std
  std.proc
    std.proc.attach
      fn (pid: i32) -> result[trace_handle, string]
      + attaches the tracer to a running process
      - returns error when the process does not exist or permission is denied
      # process
    std.proc.detach
      fn (handle: trace_handle) -> void
      + releases the traced process
      # process
    std.proc.wait_event
      fn (handle: trace_handle) -> result[raw_event, string]
      + blocks until the next syscall or signal event
      - returns error when the process has exited
      # process
  std.fmt
    std.fmt.pad_right
      fn (s: string, width: i32) -> string
      + pads a string with spaces to the given width
      + returns the original string when already at or beyond width
      # formatting

ctrace
  ctrace.start
    fn (pid: i32) -> result[trace_session, string]
    + attaches and initializes a trace session
    - returns error when attachment fails
    # lifecycle
    -> std.proc.attach
  ctrace.stop
    fn (session: trace_session) -> void
    + detaches and closes the session
    # lifecycle
    -> std.proc.detach
  ctrace.next
    fn (session: trace_session) -> result[trace_event, string]
    + returns the next decoded syscall or signal event
    - returns error when the traced process has exited
    # capture
    -> std.proc.wait_event
  ctrace.decode_syscall
    fn (raw: raw_event) -> trace_event
    + converts a raw kernel event into a named syscall with typed arguments
    # decoding
  ctrace.format_event
    fn (event: trace_event) -> string
    + returns a single-line aligned representation with name, args, and result
    # formatting
    -> std.fmt.pad_right
