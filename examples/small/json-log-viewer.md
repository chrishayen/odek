# Requirement: "a json log viewer"

A library that parses JSON log lines and renders a human-readable view. The caller handles I/O.

std
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

log_viewer
  log_viewer.parse_line
    @ (raw: string) -> result[map[string, string], string]
    + returns fields for a valid JSON log line
    - returns error for a non-JSON line
    # parsing
    -> std.json.parse_object
  log_viewer.format_line
    @ (fields: map[string, string], highlight_keys: list[string]) -> string
    + renders fields as "key=value" pairs separated by spaces
    + highlighted keys appear first in the output
    ? rendering is plain text; color codes are the caller's concern
    # formatting
