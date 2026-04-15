# Requirement: "a streaming csv parser"

Parses CSV record-by-record from an incremental byte source. Maintains a small state so records can straddle buffer boundaries.

std: (all units exist)

csv_stream
  csv_stream.new
    fn (delimiter: u8, quote: u8) -> parser_state
    + constructs a parser with the given delimiter and quote byte
    ? default usage is delimiter=',' quote='"'
    # construction
  csv_stream.feed
    fn (state: parser_state, chunk: bytes) -> tuple[parser_state, list[list[string]]]
    + returns the updated state and every record completed by this chunk
    ? incomplete trailing records are buffered in state
    # incremental_parse
  csv_stream.finish
    fn (state: parser_state) -> result[list[list[string]], string]
    + returns any final record still buffered
    - returns error when the input ends inside an unclosed quoted field
    # finalize
  csv_stream.parse_all
    fn (delimiter: u8, quote: u8, input: bytes) -> result[list[list[string]], string]
    + returns every record in a single input buffer
    - returns error on malformed quoting
    # convenience
  csv_stream.records_to_maps
    fn (records: list[list[string]]) -> result[list[map[string,string]], string]
    + treats the first record as a header and maps subsequent rows
    - returns error when a row has a different column count than the header
    - returns error when records is empty
    # shaping
