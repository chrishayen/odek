# Requirement: "a standalone event-driven TCP/IP stack suitable for bare-metal use"

The library owns no threads and performs no I/O. The caller feeds it inbound frames and drains outbound frames around a central poll.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns monotonic time in milliseconds
      # time
  std.checksum
    std.checksum.inet16
      fn (data: bytes) -> u16
      + returns the 16-bit one's complement internet checksum
      # checksum

netstack
  netstack.new_interface
    fn (mac: bytes, ipv4: u32, mtu: i32) -> iface_state
    + creates an interface with the given link and network addresses
    # construction
  netstack.inject_frame
    fn (ifc: iface_state, frame: bytes) -> result[void, string]
    + enqueues an inbound ethernet frame for the next poll
    - returns error when the frame is shorter than an ethernet header
    # ingress
  netstack.drain_frame
    fn (ifc: iface_state) -> optional[bytes]
    + returns the next outbound frame if any is queued
    # egress
  netstack.poll
    fn (ifc: iface_state) -> result[void, string]
    + processes queued frames, advances timers, runs retransmits and fires sockets' callbacks
    - returns error on an unrecoverable internal inconsistency
    # scheduling
    -> std.time.now_millis
  netstack.parse_ethernet
    fn (frame: bytes) -> result[eth_header, string]
    + returns destination, source and ethertype
    - returns error when the frame is truncated
    # parsing
  netstack.parse_ipv4
    fn (payload: bytes) -> result[ipv4_header, string]
    + returns a validated ipv4 header
    - returns error on a bad internet checksum
    # parsing
    -> std.checksum.inet16
  netstack.parse_tcp
    fn (segment: bytes, src_ip: u32, dst_ip: u32) -> result[tcp_header, string]
    + returns a validated tcp header including options
    - returns error on a bad tcp checksum
    # parsing
    -> std.checksum.inet16
  netstack.tcp_listen
    fn (ifc: iface_state, port: u16) -> result[listen_handle, string]
    + binds a listening socket on the given port
    - returns error when the port is already in use
    # sockets
  netstack.tcp_connect
    fn (ifc: iface_state, dst_ip: u32, dst_port: u16) -> result[conn_handle, string]
    + initiates an active open to the destination
    - returns error when no ephemeral port is available
    # sockets
  netstack.tcp_send
    fn (conn: conn_handle, data: bytes) -> result[i32, string]
    + appends data to the send buffer and returns the number of bytes accepted
    - returns error when the connection is closed
    # sockets
  netstack.tcp_recv
    fn (conn: conn_handle, max: i32) -> result[bytes, string]
    + returns up to max bytes from the receive buffer, possibly empty
    - returns error when the connection is closed with data remaining
    # sockets
  netstack.tcp_close
    fn (conn: conn_handle) -> result[void, string]
    + initiates an orderly shutdown of the connection
    # sockets
  netstack.step_timers
    fn (ifc: iface_state) -> void
    + advances retransmit, delayed-ack and time-wait timers based on the current clock
    # scheduling
    -> std.time.now_millis
