# Requirement: "a unified browser automation library that drives multiple headless browser engines through a single API"

The project defines a browser-agnostic facade. Each engine is represented by an opaque driver handle, and the std layer provides the transport and subprocess primitives needed to control it.

std
  std.process
    std.process.spawn
      fn (program: string, args: list[string]) -> result[process_handle, string]
      + launches a child process and returns a handle
      - returns error when the program cannot be executed
      # process
    std.process.kill
      fn (handle: process_handle) -> result[void, string]
      + terminates the child process
      # process
  std.net
    std.net.ws_connect
      fn (url: string) -> result[ws_conn, string]
      + opens a websocket connection
      # networking
    std.net.ws_send
      fn (conn: ws_conn, message: string) -> result[void, string]
      + sends a text frame
      # networking
    std.net.ws_recv
      fn (conn: ws_conn) -> result[string, string]
      + blocks until a text frame arrives and returns its payload
      # networking
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

browser
  browser.launch
    fn (engine: string, headless: bool) -> result[browser_handle, string]
    + starts an engine (one of "chromium", "webkit", "firefox") and returns a handle
    - returns error when engine is unrecognized or fails to start
    # lifecycle
    -> std.process.spawn
    -> std.net.ws_connect
  browser.close
    fn (browser: browser_handle) -> result[void, string]
    + shuts down the engine and releases resources
    # lifecycle
    -> std.process.kill
  browser.new_page
    fn (browser: browser_handle) -> result[page_handle, string]
    + opens a new tab and returns a page handle
    # navigation
  browser.goto
    fn (page: page_handle, url: string) -> result[void, string]
    + navigates the page to url and waits for load
    - returns error when navigation fails
    # navigation
    -> std.net.ws_send
    -> std.net.ws_recv
  browser.query_selector
    fn (page: page_handle, css: string) -> result[optional[element_handle], string]
    + returns the first element matching the CSS selector, or none
    # dom
  browser.click
    fn (element: element_handle) -> result[void, string]
    + dispatches a click event on the element
    - returns error when the element is detached
    # interaction
  browser.type_text
    fn (element: element_handle, text: string) -> result[void, string]
    + types text into the focused element one character at a time
    # interaction
  browser.eval_js
    fn (page: page_handle, script: string) -> result[string, string]
    + evaluates JavaScript in the page and returns the serialized result
    # scripting
    -> std.json.encode_object
    -> std.json.parse_object
  browser.screenshot
    fn (page: page_handle) -> result[bytes, string]
    + captures a PNG of the current viewport
    # capture
  browser.wait_for_selector
    fn (page: page_handle, css: string, timeout_ms: i64) -> result[element_handle, string]
    + polls until a matching element appears or timeout elapses
    - returns error on timeout
    # waiting
