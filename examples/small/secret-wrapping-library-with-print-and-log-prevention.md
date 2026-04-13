# Requirement: "a secret-wrapping library that prevents secrets from being printed or logged"

Wraps a string in an opaque container whose default string rendering is a redaction placeholder.

std: (all units exist)

secret
  secret.wrap
    @ (value: string) -> secret_state
    + returns an opaque holder for the given value
    # construction
  secret.reveal
    @ (s: secret_state) -> string
    + returns the underlying value; the only way to access it
    ? callers must audit every call site of reveal
    # access
  secret.display
    @ (s: secret_state) -> string
    + returns the placeholder string "<redacted>" regardless of contents
    + two distinct secrets produce identical display output
    # redaction
  secret.equal
    @ (a: secret_state, b: secret_state) -> bool
    + returns true when both secrets wrap the same value
    ? comparison is constant-time over the shorter input
    # comparison
