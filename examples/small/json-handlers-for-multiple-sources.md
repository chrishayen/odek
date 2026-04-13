# Requirement: "a library exposing simple handlers to read and write JSON from various sources"

Two project handlers wrap generic std I/O and JSON primitives so callers can pull JSON from a source and push it to a sink.

std
  std.io
    std.io.read_all
      @ (source: string) -> result[bytes, string]
      + reads all bytes from a named source (file path or URL)
      - returns error when the source cannot be opened
      # io
    std.io.write_all
      @ (sink: string, data: bytes) -> result[void, string]
      + writes bytes to a named sink
      - returns error when the sink cannot be opened
      # io
  std.json
    std.json.parse
      @ (raw: bytes) -> result[map[string, string], string]
      + parses JSON bytes into a string-to-string map
      - returns error on malformed JSON
      # serialization
    std.json.encode
      @ (value: map[string, string]) -> bytes
      + encodes a string-to-string map as JSON bytes
      # serialization

json_handlers
  json_handlers.read
    @ (source: string) -> result[map[string, string], string]
    + reads bytes from the source and parses them as a JSON object
    - returns error when the source is unreadable
    - returns error when the content is not valid JSON
    # read
    -> std.io.read_all
    -> std.json.parse
  json_handlers.write
    @ (sink: string, value: map[string, string]) -> result[void, string]
    + encodes the value as JSON and writes it to the sink
    - returns error when the sink cannot be written
    # write
    -> std.json.encode
    -> std.io.write_all
