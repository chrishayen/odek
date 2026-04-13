# Requirement: "a library for writing JSON to a file atomically"

Serializes a value as JSON and writes it durably by staging to a temp file and renaming. Two project runes; serialization and filesystem primitives live in std.

std
  std.json
    std.json.encode
      @ (value: json_value) -> string
      + encodes a JSON value (object, array, string, number, bool, null) to text
      # serialization
  std.fs
    std.fs.write_all
      @ (path: string, data: string) -> result[void, string]
      + writes data to path, truncating any existing file
      - returns error on permission or IO failure
      # filesystem
    std.fs.rename
      @ (from_path: string, to_path: string) -> result[void, string]
      + atomically renames a file within the same filesystem
      - returns error when source does not exist
      # filesystem
    std.fs.remove
      @ (path: string) -> result[void, string]
      + removes a file if it exists
      # filesystem

atomic_json
  atomic_json.write
    @ (path: string, value: json_value) -> result[void, string]
    + encodes the value and writes it to path atomically via a sibling temp file
    - returns error when encoding or any filesystem step fails
    ? temp file lives next to the target so rename stays within one filesystem
    # atomic_write
    -> std.json.encode
    -> std.fs.write_all
    -> std.fs.rename
    -> std.fs.remove
