# Requirement: "a lightweight data collector"

The collector ingests lines from pluggable inputs, parses them with a pluggable decoder, and buffers decoded records for a pluggable sink.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns unix time in milliseconds
      # time
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string-to-string map
      - returns error on malformed JSON or non-object root
      # serialization
    std.json.encode_object
      @ (obj: map[string,string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

collector
  collector.new
    @ (buffer_capacity: i32) -> collector_state
    + creates a collector with the given ring buffer capacity
    # construction
  collector.decode_json_line
    @ (line: string) -> result[record, string]
    + parses a JSON-formatted log line into a record with a timestamp and fields
    - returns error on malformed JSON
    - returns error when the line is missing a timestamp
    # decoding
    -> std.json.parse_object
  collector.decode_syslog_line
    @ (line: string) -> result[record, string]
    + parses an RFC3164-style syslog line into a record
    - returns error when the line lacks a header
    # decoding
    -> std.time.now_millis
  collector.ingest
    @ (state: collector_state, rec: record) -> collector_state
    + appends a decoded record to the buffer, evicting the oldest when full
    # buffering
  collector.drain
    @ (state: collector_state) -> tuple[list[record], collector_state]
    + returns every buffered record and clears the buffer
    # buffering
  collector.encode_record
    @ (rec: record) -> string
    + renders a record as a JSON line for transport to a sink
    # encoding
    -> std.json.encode_object
