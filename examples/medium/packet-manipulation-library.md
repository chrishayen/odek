# Requirement: "a packet manipulation library"

Build, parse, and serialize network packets at the byte level. Focus on a simple layered model.

std
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding
    std.encoding.hex_decode
      @ (encoded: string) -> result[bytes, string]
      - returns error on non-hex characters
      - returns error on odd length
      # encoding
  std.bytes
    std.bytes.read_u16_be
      @ (data: bytes, offset: i32) -> result[u16, string]
      + reads a big-endian u16 at offset
      - returns error when offset + 2 exceeds length
      # byte_reading
    std.bytes.write_u16_be
      @ (buf: bytes, offset: i32, value: u16) -> result[bytes, string]
      + writes a big-endian u16 at offset
      - returns error when offset + 2 exceeds length
      # byte_writing

packet
  packet.new_layer
    @ (name: string, payload: bytes) -> packet_layer
    + creates a named layer holding a raw payload
    # construction
  packet.stack
    @ (layers: list[packet_layer]) -> packet_stack
    + composes layers into an ordered stack, outermost first
    # composition
  packet.serialize
    @ (stack: packet_stack) -> bytes
    + concatenates layer payloads in order
    # serialization
  packet.parse_ethernet
    @ (raw: bytes) -> result[ethernet_frame, string]
    + extracts destination, source, and ethertype from a 14-byte header
    - returns error when raw is shorter than 14 bytes
    # parsing
    -> std.bytes.read_u16_be
  packet.build_ethernet
    @ (dst: bytes, src: bytes, ethertype: u16, payload: bytes) -> result[bytes, string]
    + constructs an ethernet frame with header and payload
    - returns error when dst or src is not 6 bytes
    # building
    -> std.bytes.write_u16_be
  packet.hexdump
    @ (raw: bytes) -> string
    + renders bytes as a hex string for inspection
    # debugging
    -> std.encoding.hex_encode
