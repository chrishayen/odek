# Requirement: "a higher-level api for building ssh servers"

A thin server layer on top of a low-level ssh handshake primitive, with handler registration and session plumbing.

std
  std.net
    std.net.tcp_listen
      @ (host: string, port: i32) -> result[listener_handle, string]
      + binds a TCP listener to the given host and port
      - returns error when the port is already in use
      # networking
    std.net.tcp_accept
      @ (lis: listener_handle) -> result[conn_handle, string]
      + blocks until a new TCP connection arrives and returns it
      # networking
  std.crypto
    std.crypto.ssh_handshake
      @ (conn: conn_handle, host_key: bytes) -> result[ssh_session, string]
      + performs the ssh transport handshake with the given host key
      - returns error on protocol mismatch or auth failure
      # cryptography
    std.crypto.load_host_key
      @ (pem: string) -> result[bytes, string]
      + parses a PEM-encoded private host key
      - returns error on malformed PEM
      # cryptography

ssh_server
  ssh_server.new
    @ (host_key_pem: string) -> result[server_state, string]
    + creates a server configured with a host key
    - returns error when the host key is invalid
    # construction
    -> std.crypto.load_host_key
  ssh_server.handle
    @ (srv: server_state, command: string, callback_id: i64) -> server_state
    + registers a handler for an exact-match command
    # routing
  ssh_server.default_handler
    @ (srv: server_state, callback_id: i64) -> server_state
    + registers a fallback handler invoked when no command matches
    # routing
  ssh_server.listen_and_serve
    @ (srv: server_state, host: string, port: i32) -> result[void, string]
    + accepts connections and dispatches sessions to registered handlers
    - returns error when the listener cannot be bound
    # serving
    -> std.net.tcp_listen
    -> std.net.tcp_accept
    -> std.crypto.ssh_handshake
  ssh_server.dispatch
    @ (srv: server_state, sess: ssh_session, command: string) -> i64
    + resolves the callback id for a command, falling back to the default
    - returns 0 when no handler is registered
    # routing
