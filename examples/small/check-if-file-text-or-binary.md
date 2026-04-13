# Requirement: "a library that determines whether a file is text or binary"

Classification uses extension hints, a BOM check, and a null-byte heuristic.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file contents
      - returns error when the path does not exist
      # filesystem

file_kind
  file_kind.classify_path
    @ (path: string) -> result[file_classification, string]
    + returns text when the file matches a known text extension, the binary path for known binary extensions, and falls through to content inspection otherwise
    # classification
    -> std.fs.read_all
  file_kind.classify_bytes
    @ (data: bytes) -> file_classification
    + returns text when the leading bytes are a known BOM
    + returns binary when a null byte appears within the first 512 bytes
    + returns text when no null byte is found in the inspection window
    ? the inspection window is capped at 512 bytes
    # classification
  file_kind.is_text_extension
    @ (extension: string) -> bool
    + returns true for extensions commonly associated with text files
    # hints
  file_kind.is_binary_extension
    @ (extension: string) -> bool
    + returns true for extensions commonly associated with binary files
    # hints
