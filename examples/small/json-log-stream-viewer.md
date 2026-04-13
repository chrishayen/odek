# Requirement: "a library for viewing structured json log streams"

Parses a stream of JSON log lines into a queryable model; rendering is the host's job.

std
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses one JSON object into a flat string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

json_log_viewer
  json_log_viewer.new
    @ () -> log_view_state
    + creates an empty log view
    # construction
  json_log_viewer.ingest_line
    @ (state: log_view_state, line: string) -> log_view_state
    + parses and appends one log line; malformed lines are kept as raw
    # ingestion
    -> std.json.parse_object
  json_log_viewer.filter
    @ (state: log_view_state, field: string, substring: string) -> list[log_entry]
    + returns entries whose named field contains substring
    - returns an empty list when no entry has that field
    # query
  json_log_viewer.columns
    @ (state: log_view_state) -> list[string]
    + returns the union of field names seen across all entries, in first-seen order
    # introspection
