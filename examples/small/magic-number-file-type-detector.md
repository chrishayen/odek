# Requirement: "a file type detector based on magic number signatures"

Two project functions: one that inspects a byte buffer and returns a type descriptor, and one that checks whether the buffer matches a specific type.

std: (all units exist)

file_type
  file_type.detect
    fn (data: bytes) -> optional[file_type_info]
    + returns a descriptor with mime type and extension for a known signature
    - returns none when no signature matches
    ? detection inspects only the first 64 bytes
    # detection
  file_type.is
    fn (data: bytes, expected_mime: string) -> bool
    + returns true when the data matches the expected mime type
    - returns false when no signature matches or the mime differs
    # check
