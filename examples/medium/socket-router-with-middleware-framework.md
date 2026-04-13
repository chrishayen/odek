# Requirement: "a low-level socket router and middleware framework"

Routes incoming socket frames to handlers by channel, with a middleware chain applied before each handler.

std
  std.net
    std.net.listen_tcp
      @ (host: string, port: u16) -> result[listener_state, string]
      + binds a TCP socket on the given host and port
      - returns error when the port is already in use
      # networking
    std.net.accept
      @ (listener: listener_state) -> result[socket_state, string]
      + blocks until a client connects and returns the new socket
      - returns error when the listener is closed
      # networking
    std.net.read_frame
      @ (sock: socket_state) -> result[bytes, string]
      + reads one length-prefixed frame from the socket
      - returns error on short read or connection close
      # networking
    std.net.write_frame
      @ (sock: socket_state, payload: bytes) -> result[void, string]
      + writes a length-prefixed frame to the socket
      # networking

socket_router
  socket_router.new
    @ () -> router_state
    + creates a router with no routes and no middleware
    # construction
  socket_router.route
    @ (state: router_state, channel: string, handler: string) -> router_state
    + registers a handler identifier for the given channel name
    + overwrites any existing handler on that channel
    # routing
  socket_router.use
    @ (state: router_state, middleware: string) -> router_state
    + appends a middleware identifier to the chain
    ? middleware runs in registration order before the handler
    # middleware
  socket_router.dispatch
    @ (state: router_state, channel: string, frame: bytes) -> result[bytes, string]
    + runs the middleware chain then the handler for the channel
    - returns error when no handler is registered for the channel
    # dispatch
  socket_router.serve
    @ (state: router_state, host: string, port: u16) -> result[void, string]
    + accepts connections and dispatches each frame by its channel prefix
    - returns error when the bind fails
    # serving
    -> std.net.listen_tcp
    -> std.net.accept
    -> std.net.read_frame
    -> std.net.write_frame
