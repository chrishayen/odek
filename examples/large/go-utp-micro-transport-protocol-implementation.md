# Requirement: "a micro transport protocol implementation providing reliable, ordered delivery over datagrams"

A reliable stream protocol layered on top of unordered datagrams: connection handshake, sequence numbers, acks, retransmission, and a congestion window.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.random
    std.random.u16
      @ () -> u16
      + returns a uniformly random 16-bit value
      # random
  std.net
    std.net.send_datagram
      @ (addr: string, payload: bytes) -> result[void, string]
      + sends a single datagram to a remote address
      - returns error on unreachable destination
      # datagram_io
    std.net.recv_datagram
      @ () -> result[tuple[string, bytes], string]
      + returns the sender address and payload for the next incoming datagram
      # datagram_io

utp
  utp.encode_packet
    @ (packet_type: u8, seq: u16, ack: u16, window: u32, payload: bytes) -> bytes
    + serializes a packet header followed by the payload
    # wire_format
  utp.decode_packet
    @ (raw: bytes) -> result[utp_packet, string]
    + parses a datagram into a structured packet
    - returns error on truncated or malformed headers
    # wire_format
  utp.open_connection
    @ (remote: string) -> result[utp_conn, string]
    + performs a three-way handshake with the remote peer
    - returns error when the peer does not reply within the handshake window
    # handshake
    -> std.random.u16
    -> std.net.send_datagram
    -> std.net.recv_datagram
  utp.accept_connection
    @ (initial: utp_packet, from_addr: string) -> result[utp_conn, string]
    + completes a handshake initiated by a remote peer
    # handshake
    -> std.net.send_datagram
  utp.send_bytes
    @ (conn: utp_conn, data: bytes) -> result[utp_conn, string]
    + breaks data into packets, enqueues them into the send window, and transmits
    + blocks when the congestion window is saturated
    # send
    -> std.net.send_datagram
    -> std.time.now_millis
  utp.recv_bytes
    @ (conn: utp_conn, max_len: i32) -> result[tuple[bytes, utp_conn], string]
    + returns the next contiguous chunk of in-order payload bytes
    + buffers out-of-order packets until the missing gap arrives
    # receive
    -> std.net.recv_datagram
  utp.process_ack
    @ (conn: utp_conn, ack_num: u16) -> utp_conn
    + removes acked packets from the send buffer
    + updates congestion window based on ack cadence
    # congestion_control
    -> std.time.now_millis
  utp.handle_timeout
    @ (conn: utp_conn) -> utp_conn
    + detects unacked packets past their deadline and retransmits them
    + halves the congestion window on loss
    # retransmission
    -> std.time.now_millis
    -> std.net.send_datagram
  utp.close_connection
    @ (conn: utp_conn) -> result[void, string]
    + sends a fin packet and waits for its ack
    # lifecycle
    -> std.net.send_datagram
    -> std.net.recv_datagram
