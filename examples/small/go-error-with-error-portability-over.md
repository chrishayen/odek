# Requirement: "a structured error library with wire-portable errors"

Errors are values with a code, message, and cause chain, and can be serialized so they survive a network hop intact.

std
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

errors
  errors.new
    @ (code: string, message: string) -> error_state
    + creates a structured error with the given code and message
    # construction
  errors.wrap
    @ (cause: error_state, code: string, message: string) -> error_state
    + wraps cause with a new outer error, preserving the chain
    # chaining
  errors.unwrap
    @ (state: error_state) -> optional[error_state]
    + returns the inner cause when present
    - returns none for a leaf error
    # chaining
  errors.encode
    @ (state: error_state) -> string
    + serializes the full chain to a portable string form
    # wire
    -> std.json.encode_object
  errors.decode
    @ (raw: string) -> result[error_state, string]
    + reconstructs an error chain from its portable form
    - returns error when the encoding is malformed
    # wire
    -> std.json.parse_object
