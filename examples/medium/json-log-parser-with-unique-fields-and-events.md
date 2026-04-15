# Requirement: "a library to parse json-like log lines and collect unique fields and events"

Reads a log stream, parses each line as a loose JSON object, and maintains running aggregates of distinct field names, distinct values per field, and event counts.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when the file cannot be opened
      # filesystem
  std.text
    std.text.split_lines
      fn (s: string) -> list[string]
      + splits on newline and drops a trailing empty segment
      # strings
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

logstats
  logstats.new
    fn () -> stats_state
    + creates an empty aggregate with zero records observed
    # construction
  logstats.ingest_line
    fn (state: stats_state, line: string) -> stats_state
    + parses the line and updates field counts, per-field distinct values, and the event counter
    ? unparseable lines increment a separate parse-error counter without aborting
    # ingestion
    -> std.json.parse_object
  logstats.ingest_stream
    fn (state: stats_state, raw: string) -> stats_state
    + folds ingest_line over every line in the stream
    # ingestion
    -> std.text.split_lines
  logstats.unique_fields
    fn (state: stats_state) -> list[string]
    + returns the distinct field names observed in insertion order
    # inspection
  logstats.unique_values
    fn (state: stats_state, field: string) -> list[string]
    + returns the distinct values observed for the named field
    + returns an empty list when the field was never observed
    # inspection
  logstats.summary
    fn (state: stats_state) -> stats_summary
    + returns totals for records ingested, parse errors, and distinct-field count
    # reporting
