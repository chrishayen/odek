# Requirement: "an IP packet decoder and analyzer library"

Decodes raw IP packets into a structured representation for analysis. Focuses on header parsing for IPv4, IPv6, TCP, UDP, and ICMP.

std
  std.bytes
    std.bytes.read_u16_be
      fn (data: bytes, offset: i32) -> result[u16, string]
      + reads a big-endian u16 at offset
      - returns error when offset is out of range
      # binary
    std.bytes.read_u32_be
      fn (data: bytes, offset: i32) -> result[u32, string]
      + reads a big-endian u32 at offset
      - returns error when offset is out of range
      # binary

packet
  packet.decode
    fn (raw: bytes) -> result[decoded_packet, string]
    + detects IPv4 vs IPv6 by version nibble and parses headers accordingly
    - returns error on truncated headers
    # decoding
    -> std.bytes.read_u16_be
    -> std.bytes.read_u32_be
  packet.decode_ipv4
    fn (raw: bytes) -> result[ipv4_header, string]
    + parses an IPv4 header including options length
    - returns error on malformed IHL
    # decoding
    -> std.bytes.read_u16_be
  packet.decode_ipv6
    fn (raw: bytes) -> result[ipv6_header, string]
    + parses a fixed 40-byte IPv6 header
    - returns error when input is shorter than 40 bytes
    # decoding
    -> std.bytes.read_u16_be
  packet.decode_tcp
    fn (raw: bytes, offset: i32) -> result[tcp_segment, string]
    + parses ports, sequence numbers, and flags from a TCP segment
    - returns error on truncated input
    # decoding
    -> std.bytes.read_u16_be
    -> std.bytes.read_u32_be
  packet.decode_udp
    fn (raw: bytes, offset: i32) -> result[udp_datagram, string]
    + parses a UDP datagram header
    - returns error on truncated input
    # decoding
    -> std.bytes.read_u16_be
  packet.summary
    fn (pkt: decoded_packet) -> string
    + returns a human-readable one-line summary including addresses and ports
    # reporting
