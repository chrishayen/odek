# Requirement: "an implementation of the address resolution protocol"

Encodes and decodes ARP frames and maintains a resolution cache mapping network addresses to hardware addresses.

std
  std.bytes
    std.bytes.read_u16_be
      fn (data: bytes, offset: i32) -> result[u16, string]
      + reads a big-endian unsigned 16-bit integer at offset
      - returns error when offset is out of range
      # byte_parsing
    std.bytes.write_u16_be
      fn (data: bytes, offset: i32, value: u16) -> bytes
      + writes a big-endian unsigned 16-bit integer at offset
      # byte_parsing
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

arp
  arp.encode_request
    fn (sender_hw: bytes, sender_ip: bytes, target_ip: bytes) -> bytes
    + encodes an ARP request frame
    # framing
    -> std.bytes.write_u16_be
  arp.encode_reply
    fn (sender_hw: bytes, sender_ip: bytes, target_hw: bytes, target_ip: bytes) -> bytes
    + encodes an ARP reply frame
    # framing
    -> std.bytes.write_u16_be
  arp.decode
    fn (frame: bytes) -> result[arp_frame, string]
    + decodes a raw ARP frame into a structured record
    - returns error when frame is shorter than the minimum header size
    - returns error when the opcode is neither request nor reply
    # parsing
    -> std.bytes.read_u16_be
  arp.new_cache
    fn (ttl_seconds: i64) -> arp_cache_state
    + returns an empty cache with the given entry lifetime
    # caching
  arp.remember
    fn (cache: arp_cache_state, ip: bytes, hw: bytes) -> arp_cache_state
    + stores or refreshes an ip-to-hardware mapping with the current timestamp
    # caching
    -> std.time.now_seconds
  arp.lookup
    fn (cache: arp_cache_state, ip: bytes) -> optional[bytes]
    + returns the hardware address for an ip when the entry has not expired
    - returns empty when the entry is missing or expired
    # caching
    -> std.time.now_seconds
