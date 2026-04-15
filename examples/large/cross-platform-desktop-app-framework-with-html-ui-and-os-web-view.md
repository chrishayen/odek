# Requirement: "a cross-platform desktop app framework with HTML UI using the OS web view"

Wraps the host OS web view into a single API for creating windows, loading content, and bridging calls between the web frontend and the native backend.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file's entire contents
      - returns error when the path does not exist
      # filesystem
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string map as JSON
      # serialization
  std.os
    std.os.detect_platform
      fn () -> string
      + returns one of "macos", "windows", "linux"
      # platform

desktop_app
  desktop_app.new
    fn (title: string, width: i32, height: i32) -> app_state
    + creates an app with a single main window of the given size
    # construction
  desktop_app.load_html
    fn (state: app_state, html: string) -> app_state
    + sets the document to the given HTML string
    # content
  desktop_app.load_url
    fn (state: app_state, url: string) -> app_state
    + navigates the embedded view to the URL
    # content
  desktop_app.load_asset
    fn (state: app_state, path: string) -> result[app_state, string]
    + reads an HTML file from disk and loads it
    - returns error when the file cannot be read
    # content
    -> std.fs.read_all
  desktop_app.bind
    fn (state: app_state, name: string, handler: string) -> app_state
    + exposes a named native handler callable from JavaScript as window[name]
    ? handler receives a JSON argument and returns a JSON result
    # bridge
  desktop_app.emit
    fn (state: app_state, event: string, payload: map[string, string]) -> app_state
    + dispatches a custom event to the frontend with a JSON payload
    # bridge
    -> std.json.encode_object
  desktop_app.on_message
    fn (state: app_state, raw: string) -> result[map[string, string], string]
    + decodes an inbound bridge message from the frontend
    - returns error when the payload is not a JSON object
    # bridge
    -> std.json.parse_object
  desktop_app.run
    fn (state: app_state) -> result[void, string]
    + opens the window and enters the platform event loop until closed
    - returns error when no web view backend is available for the current OS
    # lifecycle
    -> std.os.detect_platform
  desktop_app.close
    fn (state: app_state) -> app_state
    + requests the window to close and exits the event loop
    # lifecycle
  desktop_app.set_size
    fn (state: app_state, width: i32, height: i32) -> app_state
    + resizes the main window
    # window
