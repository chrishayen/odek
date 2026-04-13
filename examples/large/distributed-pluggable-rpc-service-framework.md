# Requirement: "a distributed pluggable rpc service framework"

Users register services; clients resolve service addresses via a pluggable registry and call methods with a pluggable codec and transport.

std
  std.net
    std.net.tcp_listen
      @ (addr: string) -> result[listener, string]
      + binds a tcp listener
      - returns error when addr is already in use
      # networking
    std.net.tcp_accept
      @ (lst: listener) -> result[conn, string]
      + accepts the next connection
      # networking
    std.net.tcp_dial
      @ (addr: string) -> result[conn, string]
      + opens a tcp connection
      # networking
    std.net.read_frame
      @ (c: conn) -> result[bytes, string]
      + reads one length-prefixed frame
      # networking
    std.net.write_frame
      @ (c: conn, data: bytes) -> result[void, string]
      + writes one length-prefixed frame
      # networking

rpcx
  rpcx.new_server
    @ () -> server_state
    + returns an empty server
    # construction
  rpcx.register_service
    @ (state: server_state, service_name: string, methods: map[string, handler_fn]) -> server_state
    + binds method handlers under the given service name
    - callers will see method not found when a method is later invoked that was not registered
    # registration
  rpcx.set_codec
    @ (state: server_state, codec: codec) -> server_state
    + installs the codec used to decode requests and encode responses
    # pluggability
  rpcx.serve
    @ (state: server_state, addr: string) -> result[server_handle, string]
    + listens on addr and dispatches incoming frames to registered handlers
    # lifecycle
    -> std.net.tcp_listen
    -> std.net.tcp_accept
    -> std.net.read_frame
    -> std.net.write_frame
  rpcx.handle_frame
    @ (state: server_state, frame: bytes) -> bytes
    + decodes a request, routes it to the matching handler, encodes the response
    + responds with a structured error when method is not found
    # dispatch
  rpcx.new_client
    @ (registry: registry_handle, codec: codec) -> client_state
    + returns a client that resolves service names through the given registry
    # client
  rpcx.call
    @ (state: client_state, service_name: string, method: string, request: bytes) -> result[bytes, string]
    + resolves an endpoint, sends a request frame, waits for the response
    - returns error when the registry returns no endpoints
    - returns error when the response is a server-side error
    # client
    -> std.net.tcp_dial
    -> std.net.write_frame
    -> std.net.read_frame
  rpcx.call_with_failover
    @ (state: client_state, service_name: string, method: string, request: bytes, attempts: i32) -> result[bytes, string]
    + calls at most attempts endpoints, moving on to the next on transport failure
    - returns the last error when every endpoint fails
    # resilience
  rpcx.register_endpoint
    @ (registry: registry_handle, service_name: string, endpoint: string) -> result[void, string]
    + records that service_name is available at endpoint
    # registry
  rpcx.unregister_endpoint
    @ (registry: registry_handle, service_name: string, endpoint: string) -> result[void, string]
    + removes an endpoint from the registry
    # registry
  rpcx.resolve_endpoints
    @ (registry: registry_handle, service_name: string) -> result[list[string], string]
    + returns all known endpoints for the service
    - returns error when service_name has no registrations
    # registry
