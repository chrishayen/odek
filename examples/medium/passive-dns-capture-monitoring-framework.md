# Requirement: "a library for passive dns capture/monitoring framework"

Observes DNS traffic from a raw packet source, decodes queries and replies, and emits records to a pluggable sink.

std
  std.net
    std.net.open_packet_source
      @ (interface_name: string) -> result[packet_source, string]
      + opens a read-only packet source on the named interface
      - returns error when the interface does not exist or lacks capture permission
      # network
    std.net.next_packet
      @ (src: packet_source) -> result[bytes, string]
      + returns the next raw packet payload
      - returns error when the source has been closed
      # network
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

dns_monitor
  dns_monitor.decode_dns
    @ (packet: bytes) -> result[dns_message, string]
    + parses a UDP DNS message into header, questions, and answers
    - returns error when the packet is shorter than the DNS header
    - returns error when name compression pointers form a cycle
    # decoding
  dns_monitor.is_query
    @ (msg: dns_message) -> bool
    + returns true when the QR bit indicates a query
    # classification
  dns_monitor.to_record
    @ (msg: dns_message, observed_at_ms: i64) -> dns_record
    + flattens a dns_message into a single record suitable for downstream sinks
    # normalization
    -> std.time.now_millis
  dns_monitor.run
    @ (interface_name: string, sink: fn(dns_record) -> void) -> result[void, string]
    + loops reading packets, decoding DNS, and calling the sink for each successful record
    - returns error when the source cannot be opened
    ? malformed packets are skipped silently; the loop continues
    # capture_loop
    -> std.net.open_packet_source
    -> std.net.next_packet
