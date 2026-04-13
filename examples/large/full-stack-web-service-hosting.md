# Requirement: "a full stack for hosting web services"

An application server that accepts HTTP requests, routes them to worker processes, and manages a worker pool.

std
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body
      - returns error on malformed request line
      # http
    std.http.serialize_response
      @ (resp: http_response) -> bytes
      + returns serialized HTTP/1.1 response bytes
      # http
  std.collections
    std.collections.queue_new
      @ () -> queue_state
      + creates an empty FIFO queue of i64 tokens
      # collections
    std.collections.queue_push
      @ (q: queue_state, value: i64) -> queue_state
      + appends to the tail
      # collections
    std.collections.queue_pop
      @ (q: queue_state) -> optional[tuple[i64, queue_state]]
      - returns none when the queue is empty
      # collections
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

app_server
  app_server.new
    @ (worker_count: i32, max_pending: i32) -> server_state
    + creates a server with the given pool size and pending limit
    - returns error when worker_count is not positive
    # construction
    -> std.collections.queue_new
  app_server.register_app
    @ (server: server_state, prefix: string, app_id: string) -> server_state
    + mounts an application at a path prefix
    # routing
  app_server.enqueue_request
    @ (server: server_state, req: http_request) -> result[tuple[i64, server_state], string]
    + assigns a request id and appends it to the pending queue
    - returns error when pending exceeds max_pending
    # queueing
    -> std.collections.queue_push
    -> std.time.now_millis
  app_server.next_ready
    @ (server: server_state) -> optional[tuple[i64, http_request, server_state]]
    + returns the next pending request bound to an available worker
    - returns none when no worker is idle or queue is empty
    # scheduling
    -> std.collections.queue_pop
  app_server.mark_done
    @ (server: server_state, request_id: i64, resp: http_response) -> result[server_state, string]
    + releases the worker that was handling this request and records the response
    - returns error when request id is unknown
    # completion
    -> std.time.now_millis
  app_server.route
    @ (server: server_state, path: string) -> optional[string]
    + returns the app_id whose prefix matches the longest path segment
    # routing
  app_server.handle_raw
    @ (server: server_state, raw: bytes) -> result[tuple[i64, server_state], string]
    + parses raw request bytes and enqueues
    - returns error on parse failure
    # io
    -> std.http.parse_request
  app_server.take_response
    @ (server: server_state, request_id: i64) -> optional[bytes]
    + returns serialized response bytes for a completed request
    # io
    -> std.http.serialize_response
  app_server.stats
    @ (server: server_state) -> server_stats
    + returns pending count, active count, completed count, and uptime ms
    # metrics
    -> std.time.now_millis
  app_server.shutdown
    @ (server: server_state) -> server_state
    + marks the server as draining so no new requests are accepted
    # lifecycle
