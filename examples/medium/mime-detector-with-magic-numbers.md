# Requirement: "a MIME type detector that identifies content by magic numbers"

Matches leading bytes of a buffer against a table of known signatures. The signature table is built-in but extensible by the caller.

std
  std.bytes
    std.bytes.starts_with
      fn (data: bytes, prefix: bytes) -> bool
      + reports whether data begins with prefix
      + returns true when prefix is empty
      # binary

mime
  mime.new_detector
    fn () -> detector_state
    + creates a detector preloaded with common signatures
    # construction
  mime.register
    fn (state: detector_state, mime_type: string, offset: i32, signature: bytes) -> detector_state
    + adds a signature that must match at the given offset
    # registration
  mime.detect
    fn (state: detector_state, data: bytes) -> string
    + returns the MIME type whose longest matching signature matches
    + returns "application/octet-stream" when nothing matches
    # detection
    -> std.bytes.starts_with
  mime.detect_extension
    fn (state: detector_state, data: bytes) -> string
    + returns a conventional file extension for the detected MIME type
    + returns "" when nothing matches
    # detection
    -> std.bytes.starts_with
