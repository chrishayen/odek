# Requirement: "an AAC bitstream parser"

Reads ADTS-framed AAC streams and exposes each frame with its decoded header fields.

std
  std.io
    std.io.bit_reader_new
      @ (data: bytes) -> bit_reader_state
      + creates a bit reader positioned at the start of data
      # io
    std.io.bit_reader_read
      @ (state: bit_reader_state, count: i32) -> result[tuple[u64, bit_reader_state], string]
      + reads count bits and advances the cursor
      - returns error when fewer than count bits remain
      # io
    std.io.bit_reader_skip
      @ (state: bit_reader_state, count: i32) -> result[bit_reader_state, string]
      + advances the cursor by count bits
      - returns error when fewer than count bits remain
      # io

aac_parser
  aac_parser.new
    @ (data: bytes) -> aac_parser_state
    + wraps data in a new parser starting at byte 0
    # construction
    -> std.io.bit_reader_new
  aac_parser.find_sync
    @ (state: aac_parser_state) -> result[aac_parser_state, string]
    + advances until the next ADTS sync word (12 bits of 1)
    - returns error when end of buffer is reached without a sync word
    # framing
  aac_parser.parse_header
    @ (state: aac_parser_state) -> result[tuple[adts_header, aac_parser_state], string]
    + decodes profile, sampling frequency index, channel config, and frame length
    - returns error on invalid sampling index or reserved fields
    # header
    -> std.io.bit_reader_read
  aac_parser.read_frame
    @ (state: aac_parser_state) -> result[tuple[aac_frame, aac_parser_state], string]
    + returns the next frame's header and payload slice and advances past it
    - returns error when the frame length exceeds the remaining buffer
    # framing
    -> std.io.bit_reader_skip
  aac_parser.iterate
    @ (data: bytes) -> result[list[aac_frame], string]
    + parses all frames sequentially until end of input
    - returns error on the first malformed frame
    # iteration
