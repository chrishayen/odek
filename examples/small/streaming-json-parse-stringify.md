# Requirement: "a streaming JSON parser and stringifier"

Parses JSON incrementally into a stream of token events and serializes a stream of values back to JSON text.

std: (all units exist)

json_stream
  json_stream.new_parser
    @ () -> parser_state
    + constructs an empty parser ready to accept input
    # construction
  json_stream.feed
    @ (state: parser_state, chunk: bytes) -> result[tuple[parser_state, list[json_event]], string]
    + returns updated state and events (start_object, key, value, start_array, end_*) produced by the chunk
    - returns error on malformed JSON
    ? incomplete tokens at chunk boundaries are buffered in state
    # incremental_parse
  json_stream.finish_parser
    @ (state: parser_state) -> result[list[json_event], string]
    + returns any trailing events
    - returns error when the input ends mid-value
    # finalize_parse
  json_stream.new_writer
    @ () -> writer_state
    + constructs an empty writer
    # construction
  json_stream.write_event
    @ (state: writer_state, event: json_event) -> result[tuple[writer_state, string], string]
    + returns the updated state and the JSON fragment produced by this event
    - returns error on an invalid event sequence (e.g. value without an open container)
    # incremental_serialize
