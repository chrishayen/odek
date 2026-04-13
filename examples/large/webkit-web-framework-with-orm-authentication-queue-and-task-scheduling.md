# Requirement: "a web framework with ORM, authentication, queue, and task scheduling"

A multi-subsystem backend toolkit: routing and middleware on top of an HTTP listener, an ORM with model bindings, password-based authentication, a job queue, and a cron-style scheduler.

std
  std.http
    std.http.listen
      @ (addr: string, handler: request_handler) -> result[void, string]
      + binds to addr and dispatches incoming requests to handler
      - returns error when the address is already in use
      # http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses an HTTP request line, headers, and body
      - returns error on malformed request
      # http
    std.http.write_response
      @ (status: i32, headers: map[string,string], body: bytes) -> bytes
      + formats a full HTTP response as bytes
      # http
  std.sql
    std.sql.connect
      @ (dsn: string) -> result[db_connection, string]
      + opens a database connection from a DSN
      - returns error on unreachable host or bad credentials
      # database
    std.sql.query
      @ (conn: db_connection, sql: string, params: list[string]) -> result[list[map[string,string]], string]
      + executes a parameterized query and returns rows as maps
      - returns error on invalid SQL
      # database
    std.sql.exec
      @ (conn: db_connection, sql: string, params: list[string]) -> result[i64, string]
      + executes a non-query statement and returns affected rows
      # database
  std.crypto
    std.crypto.hash_password
      @ (password: string) -> result[string, string]
      + hashes a password with a slow, salted algorithm
      # cryptography
    std.crypto.verify_password
      @ (password: string, hash: string) -> bool
      + returns true when the password matches the stored hash
      # cryptography
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

webkit
  webkit.new_router
    @ () -> router
    + creates an empty router
    # routing
  webkit.add_route
    @ (r: router, method: string, path: string, handler: request_handler) -> router
    + registers a handler for (method, path)
    - path patterns may contain named parameters like "/users/:id"
    # routing
  webkit.serve
    @ (r: router, addr: string) -> result[void, string]
    + starts an HTTP server dispatching to the router
    - returns error when the address is already in use
    # serving
    -> std.http.listen
    -> std.http.parse_request
    -> std.http.write_response
  webkit.define_model
    @ (table: string, columns: list[string]) -> model_schema
    + declares a model bound to a table with named columns
    # orm
  webkit.find_by_id
    @ (conn: db_connection, schema: model_schema, id: string) -> result[optional[map[string,string]], string]
    + returns a single row keyed by primary id
    - returns error on database failure
    # orm
    -> std.sql.query
  webkit.insert_row
    @ (conn: db_connection, schema: model_schema, values: map[string,string]) -> result[i64, string]
    + inserts a row and returns the new row count
    # orm
    -> std.sql.exec
  webkit.register_user
    @ (conn: db_connection, email: string, password: string) -> result[void, string]
    + hashes the password and inserts a user record
    - returns error when the email already exists
    # authentication
    -> std.crypto.hash_password
    -> std.sql.exec
  webkit.authenticate
    @ (conn: db_connection, email: string, password: string) -> result[string, string]
    + returns a session token when credentials are valid
    - returns error when the password does not match
    # authentication
    -> std.crypto.verify_password
    -> std.sql.query
  webkit.new_queue
    @ () -> job_queue
    + creates an empty in-memory job queue
    # queueing
  webkit.enqueue
    @ (q: job_queue, name: string, payload: string) -> job_queue
    + appends a job to the queue
    # queueing
  webkit.dequeue
    @ (q: job_queue) -> tuple[optional[job], job_queue]
    + pops the next job if any
    # queueing
  webkit.schedule_cron
    @ (expression: string, task_name: string) -> result[scheduled_task, string]
    + registers a task to fire on a cron expression
    - returns error when the expression is malformed
    # scheduling
  webkit.due_tasks
    @ (tasks: list[scheduled_task]) -> list[string]
    + returns the names of tasks whose next fire time has passed
    # scheduling
    -> std.time.now_seconds
