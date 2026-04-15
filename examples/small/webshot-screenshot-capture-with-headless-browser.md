# Requirement: "a library for capturing screenshots of web pages using a headless browser"

Drives a headless browser to load urls and returns encoded image bytes.

std
  std.fs
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, replacing any existing contents
      - returns error on io failure
      # storage

webshot
  webshot.new_session
    fn (browser_path: string) -> result[session_handle, string]
    + launches a headless browser process and returns a session handle
    - returns error when the browser binary is missing
    # session
  webshot.capture
    fn (session: session_handle, url: string, width: i32, height: i32) -> result[bytes, string]
    + navigates the session to url, waits for load, and returns PNG bytes
    - returns error on navigation failure or timeout
    # capture
  webshot.capture_to_file
    fn (session: session_handle, url: string, width: i32, height: i32, path: string) -> result[void, string]
    + captures a url and writes the resulting PNG to disk
    # capture
    -> std.fs.write_all
  webshot.close_session
    fn (session: session_handle) -> void
    + shuts the headless browser process down cleanly
    # lifecycle
