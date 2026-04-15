# Requirement: "an error library with stack traces and error composition"

Errors carry a message, an optional cause chain, and a captured stack trace. Formatting produces a human-readable report.

std
  std.runtime
    std.runtime.capture_stack
      fn (skip: i32) -> list[stack_frame]
      + captures the current call stack, skipping the first N frames
      # diagnostics

errorx
  errorx.new
    fn (message: string) -> rich_error
    + creates an error with the given message and a captured stack
    # construction
    -> std.runtime.capture_stack
  errorx.wrap
    fn (cause: rich_error, message: string) -> rich_error
    + wraps an existing error with additional context, preserving its stack
    ? the new error points to cause so traversal can walk the chain
    # composition
  errorx.unwrap
    fn (err: rich_error) -> optional[rich_error]
    + returns the wrapped cause, or none if this is the root error
    # composition
  errorx.chain
    fn (err: rich_error) -> list[rich_error]
    + returns the full chain from outermost to innermost
    # composition
  errorx.has_cause
    fn (err: rich_error, predicate: fn(rich_error) -> bool) -> bool
    + returns true when any error in the chain satisfies the predicate
    - returns false when no match is found
    # inspection
  errorx.format
    fn (err: rich_error) -> string
    + produces a multi-line report with every message and captured frames
    + frames are indented under each corresponding cause
    # formatting
