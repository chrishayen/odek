# Requirement: "a packet sniffer library that captures frames from an interface and parses common protocols"

Opens a capture handle, reads raw frames, and decodes Ethernet / IPv4 / TCP / UDP headers into typed records.

std
  std.net
    std.net.capture_open
      fn (iface: string, snaplen: i32) -> result[capture_handle, string]
      + opens a packet capture on the named interface
      - returns error when the interface does not exist
      # network
    std.net.capture_next
      fn (h: capture_handle) -> result[bytes, string]
      + blocks for the next frame and returns its raw bytes
      # network
    std.net.capture_close
      fn (h: capture_handle) -> void
      + releases the capture handle
      # network

sniffer
  sniffer.open
    fn (iface: string, snaplen: i32) -> result[sniffer_state, string]
    + opens a sniffer on the interface with the given snapshot length
    - returns error when the interface cannot be opened
    # capture
    -> std.net.capture_open
  sniffer.close
    fn (state: sniffer_state) -> void
    + releases the underlying capture
    # capture
    -> std.net.capture_close
  sniffer.next_frame
    fn (state: sniffer_state) -> result[bytes, string]
    + returns the next raw frame from the wire
    # capture
    -> std.net.capture_next
  sniffer.parse_ethernet
    fn (frame: bytes) -> result[tuple[ethernet_header, bytes], string]
    + parses the ethernet header and returns the payload slice
    - returns error when the frame is shorter than 14 bytes
    # parsing
  sniffer.parse_ipv4
    fn (payload: bytes) -> result[tuple[ipv4_header, bytes], string]
    + parses the IPv4 header and returns the L4 payload
    - returns error on bad version or short header
    # parsing
  sniffer.parse_tcp
    fn (payload: bytes) -> result[tcp_header, string]
    + parses the TCP header
    - returns error when the data offset is invalid
    # parsing
  sniffer.parse_udp
    fn (payload: bytes) -> result[udp_header, string]
    + parses the UDP header
    - returns error when the payload is shorter than 8 bytes
    # parsing
  sniffer.decode_frame
    fn (frame: bytes) -> result[decoded_packet, string]
    + walks ethernet, IPv4, and L4 layers into a single structured record
    - returns error when any layer fails to parse
    # parsing
