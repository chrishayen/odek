# Requirement: "a library for exposing a local service via a public tunnel"

Maintains a tunnel session: registers with an upstream relay to get a public URL, then translates incoming relay frames into local requests and pushes the responses back.

std
  std.net
    std.net.dial_tcp
      fn (host: string, port: i32) -> result[socket, string]
      + opens a TCP connection
      - returns error on connection refused or unreachable
      # networking
    std.net.send
      fn (sock: socket, data: bytes) -> result[void, string]
      + sends data on the socket
      # networking
    std.net.recv
      fn (sock: socket, max: i32) -> result[bytes, string]
      + reads up to max bytes from the socket
      # networking
  std.random
    std.random.token_hex
      fn (bytes_count: i32) -> string
      + returns a random hex token of the given byte length
      # randomness

tunnel
  tunnel.open
    fn (relay_host: string, relay_port: i32, local_port: i32) -> result[tunnel_state, string]
    + dials the relay, registers a subdomain, and returns a tunnel with the public URL
    - returns error when the relay cannot be reached
    # session
    -> std.net.dial_tcp
    -> std.random.token_hex
  tunnel.public_url
    fn (t: tunnel_state) -> string
    + returns the public URL assigned by the relay
    # query
  tunnel.next_request
    fn (t: tunnel_state) -> result[optional[tunnel_request], string]
    + reads the next incoming request frame from the relay
    - returns none when the relay has no pending requests
    # forwarding
    -> std.net.recv
  tunnel.forward_response
    fn (t: tunnel_state, req_id: string, response: bytes) -> result[void, string]
    + sends a response frame back to the relay for the given request id
    # forwarding
    -> std.net.send
  tunnel.close
    fn (t: tunnel_state) -> result[void, string]
    + unregisters the tunnel and closes the relay connection
    # session
    -> std.net.send
