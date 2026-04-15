# Requirement: "a RESTful asynchronous task queue library with topics, leases, and HTTP handlers"

Durable task queue with per-topic ordering, visibility timeouts, and HTTP endpoints for producers and consumers. Storage is pluggable.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.uuid
    std.uuid.new_v4
      fn () -> string
      + returns a random uuid string
      # identity
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on non-object root
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.http
    std.http.route
      fn (method: string, path: string, handler: fn(http_request) -> http_response) -> void
      + registers a handler for method+path
      # http_server

task_queue
  task_queue.new
    fn (store: task_store) -> queue_state
    + creates a queue backed by the given storage
    # construction
  task_queue.push
    fn (state: queue_state, topic: string, payload: string, scheduled_at_ms: i64) -> result[string, string]
    + enqueues a task and returns its id
    - returns error when topic is empty
    # producer
    -> std.uuid.new_v4
    -> std.time.now_millis
  task_queue.poll
    fn (state: queue_state, topic: string, lease_ms: i64) -> result[optional[task], string]
    + returns the next due task and grants a lease, or none when nothing is ready
    - returns error on store failure
    # consumer
    -> std.time.now_millis
  task_queue.commit
    fn (state: queue_state, task_id: string) -> result[void, string]
    + marks the task complete and removes it from the queue
    - returns error when the task id is unknown or its lease expired
    # consumer
    -> std.time.now_millis
  task_queue.release
    fn (state: queue_state, task_id: string, reason: string) -> result[void, string]
    + releases a lease so the task becomes available again
    # consumer
  task_queue.expire_leases
    fn (state: queue_state) -> i32
    + returns the number of leases that have timed out and were released
    # maintenance
    -> std.time.now_millis
  task_queue.stats
    fn (state: queue_state, topic: string) -> topic_stats
    + returns pending, leased, and completed counts for a topic
    # introspection
  task_queue.register_routes
    fn (state: queue_state, base_path: string) -> void
    + wires REST routes for push, poll, commit, and release under base_path
    # http_api
    -> std.http.route
  task_queue.handle_push
    fn (state: queue_state, request: http_request) -> http_response
    + decodes a JSON body and pushes a task, returning 201 with the id
    - returns 400 when the body is not a valid task
    # http_api
    -> std.json.parse_object
    -> std.json.encode_object
  task_queue.handle_poll
    fn (state: queue_state, request: http_request) -> http_response
    + returns 200 with a JSON task body or 204 when nothing is ready
    # http_api
    -> std.json.encode_object
