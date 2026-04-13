# Requirement: "a streaming CSV parser with a callback interface"

The caller supplies a row handler; the parser invokes it for each parsed row without buffering the whole file.

std: (all units exist)

csv_stream
  csv_stream.new
    @ (delimiter: string, has_header: bool) -> csv_parser_state
    + creates a parser with the given field delimiter
    # construction
  csv_stream.feed
    @ (state: csv_parser_state, chunk: string, on_row: fn(list[string]) -> void) -> csv_parser_state
    + consumes a chunk of input and calls on_row for each completed row
    + carries over incomplete trailing rows to the next feed
    ? quoted fields may span multiple chunks
    # parsing
  csv_stream.finish
    @ (state: csv_parser_state, on_row: fn(list[string]) -> void) -> result[void, string]
    + flushes any buffered final row
    - returns error when an unterminated quoted field remains
    # finalization
  csv_stream.header
    @ (state: csv_parser_state) -> optional[list[string]]
    + returns the header row when has_header was true and it has been parsed
    # header
