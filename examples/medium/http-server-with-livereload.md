# Requirement: "a development http-server with livereload capability"

A static file server plus a watcher that pushes reload events to connected browsers via server-sent events. Injects a tiny reload script into served HTML.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when path does not exist
      # filesystem
    std.fs.watch_tree
      fn (root: string) -> result[watch_handle, string]
      + returns a handle that emits paths on change
      - returns error when root is not a directory
      # filesystem
    std.fs.next_change
      fn (handle: watch_handle) -> optional[string]
      + returns the next changed path, none when handle closed
      # filesystem
  std.http
    std.http.serve
      fn (addr: string, handler: http_handler) -> result[server_handle, string]
      + starts listening and dispatches requests to handler
      - returns error when addr is already bound
      # networking
    std.http.write_sse_event
      fn (conn: http_conn, event: string, data: string) -> result[void, string]
      + writes a server-sent event frame to the connection
      # networking

live_server
  live_server.new
    fn (root: string, addr: string) -> result[live_server_state, string]
    + returns a server rooted at the given directory bound to addr
    - returns error when root is not a directory
    # construction
  live_server.start
    fn (state: live_server_state) -> result[live_server_state, string]
    + begins watching the root and serving files
    + opens an sse endpoint at /__livereload for browsers
    # lifecycle
    -> std.http.serve
    -> std.fs.watch_tree
  live_server.handle_request
    fn (state: live_server_state, path: string) -> http_response
    + serves matching file from root with correct content type
    + for html responses, injects a script that subscribes to /__livereload
    - returns 404 when the file does not exist under root
    # request_handling
    -> std.fs.read_all
  live_server.notify_reload
    fn (state: live_server_state, changed_path: string) -> void
    + pushes a reload event to every connected sse client
    # reload
    -> std.http.write_sse_event
  live_server.stop
    fn (state: live_server_state) -> void
    + closes the listener and watcher
    # lifecycle
