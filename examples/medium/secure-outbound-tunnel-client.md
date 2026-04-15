# Requirement: "a secure outbound tunnel client"

A library that maintains a persistent outbound control channel to a remote proxy and forwards inbound requests from that channel to local services. Network socket I/O is the host's responsibility; this library owns protocol state and routing.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.random
    std.random.uuid_v4
      fn () -> string
      + returns a random UUID v4 string
      # random
  std.encoding
    std.encoding.json_encode
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.encoding.json_decode
      fn (raw: string) -> result[map[string, string], string]
      + decodes a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

tunnel
  tunnel.new_client
    fn (token: string, remote_host: string) -> client_state
    + creates a client in disconnected state bound to a remote proxy
    # construction
  tunnel.add_ingress_rule
    fn (state: client_state, hostname: string, local_target: string) -> client_state
    + maps an incoming hostname to a local upstream URL
    # routing
  tunnel.build_register_frame
    fn (state: client_state) -> string
    + returns the JSON registration frame to send on connect
    # handshake
    -> std.encoding.json_encode
    -> std.random.uuid_v4
  tunnel.on_control_frame
    fn (state: client_state, frame: string) -> result[client_state, string]
    + processes a control frame received from the proxy (ack, config update, ping)
    - returns error when the frame is not valid JSON
    - returns error when the frame type is unknown
    # control_channel
    -> std.encoding.json_decode
    -> std.time.now_millis
  tunnel.route_request
    fn (state: client_state, hostname: string, path: string) -> result[string, string]
    + returns the local target URL for an inbound request
    - returns error when no ingress rule matches the hostname
    # routing
  tunnel.next_keepalive
    fn (state: client_state) -> tuple[optional[string], client_state]
    + returns a ping frame to send if the keepalive interval has elapsed
    # keepalive
    -> std.time.now_millis
  tunnel.mark_disconnected
    fn (state: client_state) -> client_state
    + transitions the client to disconnected and schedules a reconnect backoff
    # reconnect
    -> std.time.now_millis
  tunnel.next_reconnect_delay
    fn (state: client_state) -> tuple[i64, client_state]
    + returns the next backoff delay in milliseconds and advances the attempt counter
    # reconnect
