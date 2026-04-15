# Requirement: "a devtools driver to make web automation and scraping easy"

A driver that speaks the browser devtools protocol over a websocket, exposing page navigation, DOM querying, and input synthesis.

std
  std.json
    std.json.encode
      fn (value: json_value) -> string
      + encodes a json value as a compact string
      # serialization
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a json document
      - returns error on malformed input
      # serialization
  std.net
    std.net.ws_connect
      fn (url: string) -> result[ws_conn, string]
      + opens a websocket to the given url
      - returns error on handshake failure
      # networking
    std.net.ws_send
      fn (conn: ws_conn, frame: string) -> result[void, string]
      + sends a text frame
      # networking
    std.net.ws_recv
      fn (conn: ws_conn) -> result[string, string]
      + blocks until next text frame arrives
      - returns error when connection closed
      # networking
  std.proc
    std.proc.spawn
      fn (cmd: string, args: list[string]) -> result[proc_handle, string]
      + spawns a child process
      - returns error when binary cannot be executed
      # process

devtools_driver
  devtools_driver.launch
    fn (browser_path: string, headless: bool) -> result[driver_state, string]
    + launches a browser in remote debugging mode and attaches
    - returns error when browser binary is missing
    # lifecycle
    -> std.proc.spawn
    -> std.net.ws_connect
  devtools_driver.close
    fn (state: driver_state) -> void
    + shuts the browser down and cleans up
    # lifecycle
  devtools_driver.new_page
    fn (state: driver_state) -> result[page_handle, string]
    + creates a fresh target and returns a handle to it
    # page_management
  devtools_driver.navigate
    fn (page: page_handle, url: string) -> result[void, string]
    + navigates the page to url and waits for load event
    - returns error when url fails to load
    # navigation
    -> std.json.encode
  devtools_driver.wait_for_selector
    fn (page: page_handle, selector: string, timeout_ms: i64) -> result[void, string]
    + returns when an element matching selector is present in the dom
    - returns error when timeout elapses
    # waiting
  devtools_driver.query_text
    fn (page: page_handle, selector: string) -> result[string, string]
    + returns the text content of the first matching element
    - returns error when selector matches nothing
    # dom_query
    -> std.json.encode
    -> std.json.parse
  devtools_driver.query_all_text
    fn (page: page_handle, selector: string) -> result[list[string], string]
    + returns text content of all matching elements in document order
    # dom_query
  devtools_driver.click
    fn (page: page_handle, selector: string) -> result[void, string]
    + dispatches a click at the center of the first matching element
    - returns error when selector matches nothing
    # input
  devtools_driver.type_text
    fn (page: page_handle, selector: string, text: string) -> result[void, string]
    + focuses the element and synthesizes key events for each character
    # input
  devtools_driver.screenshot
    fn (page: page_handle) -> result[bytes, string]
    + returns a png snapshot of the current viewport
    # capture
  devtools_driver.evaluate
    fn (page: page_handle, expression: string) -> result[json_value, string]
    + evaluates a javascript expression in the page context and returns its value
    - returns error when the expression throws
    # scripting
    -> std.json.encode
    -> std.json.parse
  devtools_driver.send_command
    fn (state: driver_state, method: string, params: json_value) -> result[json_value, string]
    + sends a devtools protocol command and waits for the matching response
    - returns error when the protocol returns an error object
    ? correlates responses by monotonically increasing request id
    # protocol
    -> std.net.ws_send
    -> std.net.ws_recv
    -> std.json.encode
    -> std.json.parse
