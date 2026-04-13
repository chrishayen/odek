# Requirement: "a web framework with an asynchronous networking runtime"

Core pieces: an event loop for async I/O, an HTTP request/response pipeline, and a router. Handlers are coroutines that yield on I/O.

std
  std.io
    std.io.register_read
      @ (loop: event_loop, fd: i32, callback: read_cb) -> event_loop
      + schedules callback to run when fd is readable
      # async_io
    std.io.register_write
      @ (loop: event_loop, fd: i32, callback: write_cb) -> event_loop
      + schedules callback to run when fd is writable
      # async_io
    std.io.run_once
      @ (loop: event_loop, timeout_ms: i32) -> event_loop
      + blocks until one event fires or the timeout elapses, then dispatches it
      # async_io
  std.net
    std.net.listen_tcp
      @ (host: string, port: i32) -> result[i32, string]
      + returns a listening socket file descriptor
      - returns error on bind failure
      # networking
    std.net.accept
      @ (listen_fd: i32) -> result[i32, string]
      + returns a new connected socket file descriptor
      - returns error when no pending connection is ready
      # networking
    std.net.read
      @ (fd: i32, max: i32) -> result[bytes, string]
      + reads up to max bytes without blocking
      # networking
    std.net.write
      @ (fd: i32, data: bytes) -> result[i32, string]
      + writes bytes without blocking; returns the number written
      # networking
  std.http
    std.http.parse_request
      @ (raw: string) -> result[http_request, string]
      + parses an HTTP/1.1 request
      - returns error on invalid request lines
      # http
    std.http.render_response
      @ (status: i32, headers: map[string, string], body: string) -> string
      + returns a wire-format HTTP response
      # http

web_framework
  web_framework.new_loop
    @ () -> event_loop
    + returns an empty event loop
    # construction
  web_framework.new_app
    @ () -> app_state
    + returns an empty application
    # construction
  web_framework.route
    @ (app: app_state, method: string, pattern: string, handler: async_handler) -> app_state
    + registers a handler for the method and pattern
    + supports placeholders like "/users/{id}"
    # routing
  web_framework.match_route
    @ (app: app_state, method: string, path: string) -> optional[route_match]
    + returns the matching handler and extracted parameters
    - returns none when no route matches
    # routing
  web_framework.serve
    @ (app: app_state, host: string, port: i32) -> result[event_loop, string]
    + binds the socket and registers the accept callback
    - returns error on bind failure
    # serving
    -> std.net.listen_tcp
    -> std.io.register_read
  web_framework.handle_connection
    @ (app: app_state, fd: i32, buffer: bytes) -> connection_state
    + advances a connection state machine as bytes arrive
    # serving
    -> std.http.parse_request
    -> std.http.render_response
  web_framework.gather
    @ (tasks: list[task]) -> list[task_result]
    + awaits all tasks and returns their results in order
    # concurrency
  web_framework.run_forever
    @ (loop: event_loop) -> void
    + runs the event loop until no callbacks remain
    # concurrency
    -> std.io.run_once
