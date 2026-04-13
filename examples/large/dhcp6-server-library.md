# Requirement: "a DHCPv6 server library"

Handles message parsing, lease state, and option handling. Network I/O is injected by the caller; this library only transforms messages and manages state.

std
  std.net
    std.net.parse_ipv6
      @ (value: string) -> result[bytes, string]
      + parses a textual IPv6 address into 16 bytes
      - returns error on malformed addresses
      # networking
    std.net.format_ipv6
      @ (addr: bytes) -> string
      + renders 16 bytes as a canonical IPv6 string
      # networking
  std.bytes
    std.bytes.read_u16_be
      @ (data: bytes, offset: i32) -> result[u16, string]
      + reads a big-endian u16 at offset
      - returns error when offset is out of range
      # binary
    std.bytes.read_u32_be
      @ (data: bytes, offset: i32) -> result[u32, string]
      + reads a big-endian u32 at offset
      - returns error when offset is out of range
      # binary
    std.bytes.write_u16_be
      @ (buf: bytes, value: u16) -> bytes
      + appends value as big-endian u16
      # binary
    std.bytes.write_u32_be
      @ (buf: bytes, value: u32) -> bytes
      + appends value as big-endian u32
      # binary
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

dhcp6
  dhcp6.new_server
    @ (pool_start: string, pool_end: string, lease_seconds: i64) -> result[server_state, string]
    + creates a server with an address pool and default lease time
    - returns error on invalid pool addresses
    # construction
    -> std.net.parse_ipv6
  dhcp6.parse_message
    @ (raw: bytes) -> result[dhcp6_message, string]
    + parses a client message including its transaction id and options
    - returns error on truncated or malformed input
    # parsing
    -> std.bytes.read_u16_be
    -> std.bytes.read_u32_be
  dhcp6.encode_message
    @ (msg: dhcp6_message) -> bytes
    + serializes a server response message
    # serialization
    -> std.bytes.write_u16_be
    -> std.bytes.write_u32_be
  dhcp6.handle_solicit
    @ (state: server_state, msg: dhcp6_message) -> tuple[server_state, dhcp6_message]
    + reserves a candidate address and returns an ADVERTISE response
    # protocol
  dhcp6.handle_request
    @ (state: server_state, msg: dhcp6_message) -> tuple[server_state, dhcp6_message]
    + binds the requested address to the client DUID and returns a REPLY
    # protocol
    -> std.time.now_seconds
  dhcp6.handle_release
    @ (state: server_state, msg: dhcp6_message) -> tuple[server_state, dhcp6_message]
    + frees the bound address and returns a REPLY
    # protocol
  dhcp6.handle_renew
    @ (state: server_state, msg: dhcp6_message) -> tuple[server_state, dhcp6_message]
    + extends the lease on the bound address
    # protocol
    -> std.time.now_seconds
  dhcp6.expire_leases
    @ (state: server_state) -> server_state
    + reclaims all leases whose expiry has passed
    # gc
    -> std.time.now_seconds
  dhcp6.lookup_lease
    @ (state: server_state, duid: bytes) -> optional[lease]
    + returns the active lease for a client DUID if any
    # lookup
