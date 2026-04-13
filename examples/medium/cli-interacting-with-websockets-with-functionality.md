# Requirement: "a library for bidirectional WebSocket and TCP streaming between pluggable endpoints"

An endpoint is an abstract byte stream. The library opens endpoints of several kinds, joins two streams in both directions, and surfaces lifecycle events.

std
  std.net
    std.net.dial_tcp
      @ (host: string, port: i32) -> result[stream, string]
      + returns a connected TCP stream
      - returns error when the connection is refused
      # networking
    std.net.listen_tcp
      @ (host: string, port: i32) -> result[listener, string]
      + binds a TCP listener
      - returns error when the port is already bound
      # networking

streamcat
  streamcat.open_ws
    @ (url: string) -> result[stream, string]
    + performs a WebSocket handshake and returns a message-framed stream
    - returns error when the server rejects the handshake
    # endpoints
  streamcat.open_tcp
    @ (host: string, port: i32) -> result[stream, string]
    + returns a connected TCP stream
    # endpoints
    -> std.net.dial_tcp
  streamcat.open_listen
    @ (host: string, port: i32) -> result[stream, string]
    + waits for a single TCP connection on host:port and returns it as a stream
    # endpoints
    -> std.net.listen_tcp
  streamcat.open_stdio
    @ () -> stream
    + returns a stream that reads from standard input and writes to standard output
    # endpoints
  streamcat.join
    @ (a: stream, b: stream) -> result[join_stats, string]
    + copies bytes in both directions until either side ends, returning byte counts
    - returns error when either stream fails mid-transfer
    # bridging
  streamcat.parse_endpoint
    @ (spec: string) -> result[endpoint_spec, string]
    + parses an endpoint spec like "ws://...", "tcp-listen:8080", or "-"
    - returns error on unrecognized schemes
    # parsing
  streamcat.open
    @ (spec: endpoint_spec) -> result[stream, string]
    + opens a stream matching the parsed endpoint spec
    # endpoints
    -> streamcat.open_ws
    -> streamcat.open_tcp
    -> streamcat.open_listen
    -> streamcat.open_stdio
