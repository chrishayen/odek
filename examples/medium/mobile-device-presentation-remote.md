# Requirement: "a library that lets a mobile device act as a presentation remote"

Pairs a phone with a desktop over a local network so taps on the phone advance or reverse slides on the desktop.

std
  std.net
    std.net.tcp_listen
      fn (port: u16) -> result[listener_state, string]
      + binds a TCP listener on the given port
      - returns error when the port is already in use
      # network
    std.net.tcp_accept
      fn (l: listener_state) -> result[conn_state, string]
      + blocks until a client connects
      # network
    std.net.tcp_read
      fn (c: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes from the connection
      # network
    std.net.tcp_write
      fn (c: conn_state, data: bytes) -> result[void, string]
      + writes bytes to the connection
      # network
  std.crypto
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography

remote
  remote.start_host
    fn (port: u16) -> result[host_state, string]
    + starts the desktop side listening for a single phone pairing
    # host
    -> std.net.tcp_listen
  remote.generate_pairing_code
    fn (host: host_state) -> string
    + returns a short numeric code for the user to enter on the phone
    ? derived from random bytes so codes are non-guessable
    # pairing
    -> std.crypto.random_bytes
  remote.accept_pairing
    fn (host: host_state, code: string) -> result[session_state, string]
    + waits for a phone to connect and confirms the pairing code
    - returns error when the connecting phone sends a mismatched code
    # pairing
    -> std.net.tcp_accept
    -> std.net.tcp_read
  remote.next_event
    fn (session: session_state) -> result[remote_event, string]
    + blocks until the phone sends a tap event and decodes it as NEXT, PREV, or START
    - returns error when the session is closed
    # events
    -> std.net.tcp_read
  remote.send_ack
    fn (session: session_state, event: remote_event) -> result[void, string]
    + acknowledges an event back to the phone so it can vibrate on receipt
    # events
    -> std.net.tcp_write
