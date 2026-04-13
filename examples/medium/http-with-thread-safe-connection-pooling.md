# Requirement: "an http client with a thread-safe connection pool and file upload support"

A connection pool indexed by host plus a request-sender that borrows from it, with multipart form support.

std
  std.net
    std.net.dial
      @ (host: string, port: i32) -> result[connection, string]
      + opens a tcp connection to host:port
      - returns error on connection failure
      # networking
    std.net.close
      @ (c: connection) -> void
      + closes the connection
      # networking
  std.sync
    std.sync.mutex_lock
      @ (m: mutex) -> void
      + acquires the mutex, blocking until available
      # concurrency
    std.sync.mutex_unlock
      @ (m: mutex) -> void
      + releases the mutex
      # concurrency

http_pool
  http_pool.new_pool
    @ (max_per_host: i32) -> pool_state
    + returns an empty pool with the per-host cap
    # construction
  http_pool.acquire
    @ (state: pool_state, host: string, port: i32) -> result[tuple[connection, pool_state], string]
    + returns a reused idle connection when one exists for host:port
    + dials a new connection when none are idle and the cap is not yet reached
    - blocks the caller when the cap is reached and no idle connections exist
    # pooling
    -> std.sync.mutex_lock
    -> std.sync.mutex_unlock
    -> std.net.dial
  http_pool.release
    @ (state: pool_state, host: string, port: i32, conn: connection) -> pool_state
    + returns the connection to the idle list for reuse
    - closes the connection when the idle list is full
    # pooling
    -> std.net.close
  http_pool.post_multipart
    @ (state: pool_state, url: string, fields: map[string,string], files: list[upload_file]) -> result[http_response, string]
    + sends a multipart/form-data POST, borrowing a pooled connection
    + returns the response after releasing the connection back to the pool
    - returns error when the url cannot be parsed
    # requests
  http_pool.encode_multipart
    @ (boundary: string, fields: map[string,string], files: list[upload_file]) -> bytes
    + produces the multipart body with the given boundary
    + orders fields before files
    # encoding
