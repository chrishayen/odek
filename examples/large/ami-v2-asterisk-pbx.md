# Requirement: "a client library for a telephone-exchange management and event protocol"

Speaks a text-based line protocol over TCP: login, send actions, read responses, and dispatch asynchronous events.

std
  std.net
    std.net.tcp_connect
      @ (host: string, port: u16) -> result[conn, string]
      + opens a TCP connection to host:port
      - returns error when the host is unreachable
      # networking
    std.net.read_line
      @ (c: conn) -> result[string, string]
      + reads a single line terminated by CRLF and returns it without the terminator
      - returns error on close before a line is complete
      # networking
    std.net.write_bytes
      @ (c: conn, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      - returns error on broken pipe
      # networking
  std.hash
    std.hash.md5_hex
      @ (data: bytes) -> string
      + returns the lowercase hex MD5 digest of data
      # hashing

pbx_client
  pbx_client.connect
    @ (host: string, port: u16) -> result[pbx_conn, string]
    + opens a protocol connection and reads the server banner
    - returns error when the banner is not recognized
    # connection
    -> std.net.tcp_connect
    -> std.net.read_line
  pbx_client.login
    @ (c: pbx_conn, username: string, secret: string) -> result[pbx_conn, string]
    + authenticates with a challenge/response using an MD5 digest of the secret and the server salt
    - returns error when credentials are rejected
    # authentication
    -> std.hash.md5_hex
    -> std.net.write_bytes
    -> std.net.read_line
  pbx_client.encode_message
    @ (fields: list[tuple[string, string]]) -> bytes
    + encodes action fields as "Key: Value" lines terminated by a blank line
    # protocol
  pbx_client.read_message
    @ (c: pbx_conn) -> result[list[tuple[string, string]], string]
    + reads a message from the connection until a blank line and returns its fields in order
    - returns error on connection close
    # protocol
    -> std.net.read_line
  pbx_client.send_action
    @ (c: pbx_conn, action: string, params: list[tuple[string, string]]) -> result[list[tuple[string, string]], string]
    + sends an action with an auto-assigned ActionID and returns the correlated response
    - returns error when the response reports failure
    # action_dispatch
    -> std.net.write_bytes
  pbx_client.subscribe_events
    @ (c: pbx_conn, handler: fn(list[tuple[string, string]]) -> void) -> result[void, string]
    + reads messages in a loop and dispatches event messages (non-response) to handler until the connection closes
    - returns error on read failure
    # event_loop
  pbx_client.logoff
    @ (c: pbx_conn) -> result[void, string]
    + sends a logoff action and closes the connection
    - returns error when the logoff is not acknowledged
    # teardown
