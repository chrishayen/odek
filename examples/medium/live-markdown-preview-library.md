# Requirement: "a live markdown preview library"

Watches a markdown file, re-renders on change, and pushes HTML to connected viewers over a websocket.

std
  std.fs
    std.fs.read_file
      fn (path: string) -> result[bytes, string]
      + reads file contents
      - returns error when file is missing
      # filesystem
    std.fs.watch
      fn (path: string) -> result[watcher_state, string]
      + starts watching a path for modifications
      - returns error when path cannot be watched
      # filesystem
    std.fs.poll_event
      fn (w: watcher_state) -> optional[string]
      + returns the path of the next changed file, or none when idle
      # filesystem
  std.net
    std.net.ws_listen
      fn (port: i32) -> result[ws_server, string]
      + starts a websocket server on the given port
      - returns error when the port is in use
      # networking
    std.net.ws_broadcast
      fn (srv: ws_server, payload: bytes) -> result[void, string]
      + sends a message to every connected client
      # networking

md_preview
  md_preview.render
    fn (markdown: string) -> string
    + converts markdown source to an HTML fragment
    ? supports headings, emphasis, code fences, lists, links
    # rendering
  md_preview.start
    fn (file_path: string, port: i32) -> result[preview_state, string]
    + reads the file, starts the websocket server, and begins watching
    - returns error when the file or port are unavailable
    # lifecycle
    -> std.fs.read_file
    -> std.fs.watch
    -> std.net.ws_listen
  md_preview.tick
    fn (state: preview_state) -> result[bool, string]
    + processes one pending file change, re-rendering and broadcasting when needed; returns true when an update was sent
    # loop_step
    -> std.fs.poll_event
    -> std.fs.read_file
    -> std.net.ws_broadcast
  md_preview.stop
    fn (state: preview_state) -> result[void, string]
    + stops the server and file watcher
    # lifecycle
