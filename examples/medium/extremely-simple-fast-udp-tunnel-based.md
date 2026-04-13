# Requirement: "a reliable tunnel over an unreliable datagram protocol"

Provides a stream-oriented session layered on top of a datagram transport, with sequencing, acks, and retransmission.

std
  std.udp
    std.udp.bind
      @ (addr: string) -> result[udp_socket, string]
      + binds to a local address
      - returns error on bind failure
      # transport
    std.udp.send_to
      @ (sock: udp_socket, peer: string, data: bytes) -> result[void, string]
      + sends a datagram to peer
      # transport
    std.udp.recv_from
      @ (sock: udp_socket) -> result[tuple[string, bytes], string]
      + blocks until a datagram arrives and returns (peer, payload)
      # transport
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

tunnel
  tunnel.new_session
    @ (peer: string, mtu: i32) -> session_state
    + creates a new session targeting peer with the given segment size
    # construction
  tunnel.write
    @ (state: session_state, data: bytes) -> session_state
    + enqueues bytes for transmission, splitting into numbered segments
    # send
  tunnel.on_datagram
    @ (state: session_state, payload: bytes) -> tuple[session_state, bytes]
    + processes a received datagram, returning any newly reassembled payload
    - returns empty bytes when the datagram is a duplicate or ack only
    # receive
  tunnel.tick
    @ (state: session_state) -> tuple[session_state, list[bytes]]
    + returns segments whose retransmit timer has expired
    ? retransmit timeout doubles per attempt up to a cap
    # retransmission
    -> std.time.now_millis
  tunnel.run
    @ (sock: udp_socket, state: session_state) -> result[void, string]
    + drives send, receive, and tick until the session closes
    # event_loop
    -> std.udp.send_to
    -> std.udp.recv_from
