# Requirement: "a simple error wrapper"

Wraps an underlying error with a message and allows walking the chain.

std: (all units exist)

errwrap
  errwrap.wrap
    fn (message: string, cause: optional[wrapped_error]) -> wrapped_error
    + attaches a message to an optional underlying cause
    # wrapping
  errwrap.message
    fn (err: wrapped_error) -> string
    + returns the top-level message
    # inspection
  errwrap.cause
    fn (err: wrapped_error) -> optional[wrapped_error]
    + returns the wrapped cause, if any
    # inspection
  errwrap.format_chain
    fn (err: wrapped_error) -> string
    + returns messages joined by ": " from outermost to innermost
    - returns just the top message when there is no cause
    # formatting
