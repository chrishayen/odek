# Requirement: "an error handling package"

Structured errors with a code, message, and optional cause chain.

std: (all units exist)

errs
  errs.new
    fn (code: string, message: string) -> error_value
    + creates a leaf error with the given code and message
    # construction
  errs.wrap
    fn (cause: error_value, code: string, message: string) -> error_value
    + creates an error that wraps cause
    # wrapping
  errs.code
    fn (err: error_value) -> string
    + returns the topmost error's code
    # inspection
  errs.format_chain
    fn (err: error_value) -> string
    + returns a multi-line rendering of err followed by each cause
    # formatting
  errs.has_code
    fn (err: error_value, code: string) -> bool
    + returns true when err or any cause has the given code
    # inspection
