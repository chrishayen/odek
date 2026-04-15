# Requirement: "a TCP connection hijacker"

Crafts TCP segments that inject data into an existing connection given its 4-tuple and sequence numbers. Packet capture and raw socket IO are host-level; the library produces and parses segments.

std
  std.net
    std.net.send_raw
      fn (packet: bytes) -> result[void, string]
      + sends a raw IP packet via a raw socket
      - returns error when the caller lacks permission
      # network
    std.net.capture_next
      fn (filter: string) -> result[bytes, string]
      + returns the next packet matching the bpf filter
      - returns error when the capture is closed
      # network
  std.hash
    std.hash.ones_complement_sum
      fn (data: bytes) -> u16
      + returns the 16-bit one's complement sum used by IP/TCP checksums
      # hashing

hijack
  hijack.connection
    fn (src_ip: string, src_port: u16, dst_ip: string, dst_port: u16) -> conn_info
    + builds a connection descriptor for the 4-tuple
    # construction
  hijack.observe
    fn (conn: conn_info) -> result[conn_state, string]
    + captures one packet to learn current sequence and ack numbers
    - returns error when no matching packet arrives
    # observation
    -> std.net.capture_next
  hijack.build_segment
    fn (conn: conn_info, state: conn_state, payload: bytes, flags: u8) -> bytes
    + assembles an IP+TCP packet with a correct checksum
    # assembly
    -> std.hash.ones_complement_sum
  hijack.inject
    fn (conn: conn_info, state: conn_state, payload: bytes) -> result[conn_state, string]
    + sends a payload segment and returns the advanced state
    - returns error when the send fails
    # injection
    -> std.net.send_raw
  hijack.send_reset
    fn (conn: conn_info, state: conn_state) -> result[void, string]
    + sends a RST to tear down the connection
    # teardown
    -> std.net.send_raw
  hijack.parse_segment
    fn (packet: bytes) -> result[segment_view, string]
    + decodes an IP+TCP packet into source/dest/seq/ack/flags/payload
    - returns error when the packet is too short or not TCP over IPv4
    # parsing
