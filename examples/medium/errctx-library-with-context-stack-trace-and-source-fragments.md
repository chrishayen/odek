# Requirement: "an error handling library with context, stack trace, and source fragments"

A structured error type that carries a message, wrapped cause, key/value context, a captured call stack, and source fragments read from files referenced by stack frames.

std
  std.runtime
    std.runtime.capture_stack
      @ (skip: i32) -> list[stack_frame]
      + returns the current call stack, skipping the top `skip` frames
      # runtime
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the full contents of a file
      - returns error when the file cannot be opened
      # filesystem

errctx
  errctx.new
    @ (message: string) -> error_value
    + creates an error with message and a freshly captured stack
    -> std.runtime.capture_stack
    # construction
  errctx.wrap
    @ (cause: error_value, message: string) -> error_value
    + returns a new error whose cause is `cause`, with its own stack
    -> std.runtime.capture_stack
    # wrapping
  errctx.with
    @ (err: error_value, key: string, value: string) -> error_value
    + returns a copy of the error with an additional context key/value
    # context
  errctx.unwrap
    @ (err: error_value) -> optional[error_value]
    + returns the wrapped cause if any
    - returns none at the root of the chain
    # inspection
  errctx.stack
    @ (err: error_value) -> list[stack_frame]
    + returns the stack captured at the error's origin
    # inspection
  errctx.source_fragments
    @ (err: error_value, context_lines: i32) -> list[source_fragment]
    + returns surrounding source lines for each frame whose file is readable
    ? frames whose files cannot be read are skipped silently
    -> std.fs.read_all
    # inspection
  errctx.format
    @ (err: error_value) -> string
    + returns a multiline rendering with message, context pairs, and stack
    # rendering
