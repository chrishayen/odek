# Requirement: "a low- and high-level HTTP/2 server library"

Exposes a frame-level protocol state machine as the low-level layer and a request/response router as the high-level layer.

std
  std.bytes
    std.bytes.concat
      @ (parts: list[bytes]) -> bytes
      + concatenates byte sequences in order
      # bytes
  std.io
    std.io.read_u24_be
      @ (data: bytes, offset: i32) -> result[i32, string]
      + reads a 24-bit big-endian unsigned integer at offset
      - returns error when offset + 3 exceeds length
      # binary
    std.io.read_u32_be
      @ (data: bytes, offset: i32) -> result[i32, string]
      + reads a 32-bit big-endian unsigned integer at offset
      - returns error when offset + 4 exceeds length
      # binary

http2
  http2.new_connection
    @ (is_server: bool) -> conn_state
    + creates a fresh connection state with default settings and empty stream table
    # construction
  http2.parse_frame
    @ (data: bytes, offset: i32) -> result[tuple[frame, i32], string]
    + parses one frame starting at offset and returns (frame, next_offset)
    - returns error on truncated input or unknown frame type
    # low_level
    -> std.io.read_u24_be
    -> std.io.read_u32_be
  http2.encode_frame
    @ (f: frame) -> bytes
    + serializes a frame to its wire representation
    # low_level
    -> std.bytes.concat
  http2.apply_frame
    @ (state: conn_state, f: frame) -> result[tuple[conn_state, list[frame]], string]
    + updates connection state and returns any frames that must be sent in response
    - returns error on protocol violation (e.g. DATA on idle stream)
    # low_level
  http2.open_stream
    @ (state: conn_state, headers: map[string, string]) -> tuple[i32, conn_state]
    + allocates a client-side stream id and records its pending headers
    # low_level
  http2.decode_headers
    @ (data: bytes) -> result[map[string, string], string]
    + decodes a header block into a name/value map
    - returns error on malformed input
    # header_coding
  http2.encode_headers
    @ (headers: map[string, string]) -> bytes
    + encodes headers into a wire header block
    # header_coding
  http2.new_router
    @ () -> router_state
    + creates an empty high-level router
    # high_level
  http2.route
    @ (router: router_state, method: string, path: string, handler_name: string) -> router_state
    + registers a handler name for the given method/path pair
    # high_level
  http2.dispatch
    @ (router: router_state, method: string, path: string) -> optional[string]
    + returns the handler name for a request
    - returns none when no route matches
    # high_level
  http2.build_response
    @ (status: i32, headers: map[string, string], body: bytes) -> response
    + assembles a response value suitable for serialization via encode_frame
    # high_level
