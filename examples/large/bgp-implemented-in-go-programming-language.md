# Requirement: "a BGP-4 protocol implementation"

A Border Gateway Protocol speaker: message codec, FSM per peer, and a minimal RIB with best-path selection.

std
  std.net
    std.net.tcp_connect
      @ (host: string, port: u16) -> result[conn_state, string]
      + opens a TCP connection to the remote peer
      - returns error on unreachable host or refused connection
      # networking
    std.net.tcp_read
      @ (conn: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes from the connection
      - returns error when the connection is closed
      # networking
    std.net.tcp_write
      @ (conn: conn_state, data: bytes) -> result[i32, string]
      + writes all bytes and returns the count written
      - returns error when the connection is closed
      # networking
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns the current unix time in seconds
      # time
  std.encoding
    std.encoding.read_u16_be
      @ (data: bytes, offset: i32) -> result[u16, string]
      + reads a big-endian u16 at offset
      - returns error when offset + 2 exceeds length
      # encoding
    std.encoding.write_u16_be
      @ (value: u16) -> bytes
      + returns a 2-byte big-endian representation
      # encoding
    std.encoding.read_u32_be
      @ (data: bytes, offset: i32) -> result[u32, string]
      + reads a big-endian u32 at offset
      - returns error when offset + 4 exceeds length
      # encoding

bgp
  bgp.encode_open
    @ (as_number: u16, hold_time: u16, router_id: u32) -> bytes
    + encodes a BGP OPEN message with the 19-byte header and body
    # codec
    -> std.encoding.write_u16_be
  bgp.encode_keepalive
    @ () -> bytes
    + encodes a 19-byte KEEPALIVE message with only the header
    # codec
  bgp.encode_update
    @ (withdrawn: list[prefix], announced: list[prefix], nexthop: u32, as_path: list[u16]) -> bytes
    + encodes a BGP UPDATE with withdrawn routes, path attributes, and NLRI
    # codec
    -> std.encoding.write_u16_be
  bgp.decode_message
    @ (data: bytes) -> result[bgp_message, string]
    + parses a complete BGP message (OPEN, UPDATE, NOTIFICATION, KEEPALIVE)
    - returns error when the 16-byte marker is not all ones
    - returns error when the declared length exceeds the buffer
    # codec
    -> std.encoding.read_u16_be
    -> std.encoding.read_u32_be
  bgp.fsm_new
    @ (peer_as: u16, local_as: u16, hold_time: u16) -> fsm_state
    + creates a peer FSM in the Idle state
    # session
  bgp.fsm_transition
    @ (state: fsm_state, event: fsm_event) -> fsm_state
    + advances the FSM (Idle -> Connect -> OpenSent -> OpenConfirm -> Established)
    - moves to Idle on HoldTimer expiry or NOTIFICATION
    # session
    -> std.time.now_seconds
  bgp.rib_new
    @ () -> rib_state
    + creates an empty routing information base
    # routing_table
  bgp.rib_apply_update
    @ (rib: rib_state, peer_id: u32, msg: bgp_message) -> rib_state
    + installs announced prefixes from the update keyed by (prefix, peer)
    + removes withdrawn prefixes for the peer
    # routing_table
  bgp.rib_best_path
    @ (rib: rib_state, p: prefix) -> optional[route]
    + selects the best route for a prefix using shortest AS_PATH then lowest peer id
    - returns none when no route exists for the prefix
    # path_selection
