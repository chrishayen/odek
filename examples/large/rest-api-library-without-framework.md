# Requirement: "a rest api library without a framework"

A minimal HTTP server with route registration, method matching, JSON body handling, and response helpers. Enough structure to build a REST API from scratch.

std
  std.net
    std.net.tcp_listen
      @ (host: string, port: i32) -> result[listener_state, string]
      + binds a TCP listener on host:port
      - returns error when the port is in use
      # networking
    std.net.tcp_accept
      @ (listener: listener_state) -> result[conn_state, string]
      + blocks until a client connects and returns a connection handle
      # networking
    std.net.conn_read
      @ (conn: conn_state, max: i32) -> result[bytes, string]
      + reads up to max bytes from the connection
      # networking
    std.net.conn_write
      @ (conn: conn_state, data: bytes) -> result[void, string]
      + writes the full buffer to the connection
      # networking
    std.net.conn_close
      @ (conn: conn_state) -> void
      + closes the connection
      # networking
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a dynamic value
      - returns error on malformed JSON
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + serializes a dynamic JSON value to a string
      # serialization

rest_api
  rest_api.new_router
    @ () -> router_state
    + creates an empty router with no routes
    # construction
  rest_api.add_route
    @ (router: router_state, method: string, path: string, handler: handler_fn) -> router_state
    + registers a handler for (method, path)
    ? path may contain ":name" segments which become path parameters
    # routing
  rest_api.parse_request
    @ (raw: bytes) -> result[http_request, string]
    + parses an HTTP/1.1 request line, headers, and body from raw bytes
    - returns error when the request line is malformed
    - returns error when Content-Length exceeds the body length
    # http_parsing
  rest_api.match_route
    @ (router: router_state, method: string, path: string) -> result[route_match, string]
    + returns the handler and extracted path params for the request
    - returns error when no route matches
    # routing
  rest_api.build_response
    @ (status: i32, headers: map[string,string], body: bytes) -> bytes
    + serializes a response into HTTP/1.1 wire format
    + sets Content-Length automatically from body length
    # http_serialization
  rest_api.json_response
    @ (status: i32, value: json_value) -> http_response
    + builds a response with Content-Type application/json and encoded body
    # response_helpers
    -> std.json.encode
  rest_api.read_json_body
    @ (req: http_request) -> result[json_value, string]
    + parses the request body as JSON
    - returns error when body is empty or invalid
    # request_helpers
    -> std.json.parse
  rest_api.serve
    @ (router: router_state, host: string, port: i32) -> result[void, string]
    + accepts connections and dispatches each request through the router
    - returns error when the listener cannot bind
    # server_loop
    -> std.net.tcp_listen
    -> std.net.tcp_accept
    -> std.net.conn_read
    -> std.net.conn_write
    -> std.net.conn_close
