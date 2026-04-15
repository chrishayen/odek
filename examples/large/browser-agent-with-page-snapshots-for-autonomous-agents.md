# Requirement: "a browser automation library that exposes pages as structured snapshots for autonomous agents"

Drives a browser session and returns each page as a compact structured snapshot of interactive elements keyed by stable ids so an agent can act on them by id.

std
  std.net
    std.net.http_get
      fn (url: string) -> result[string, string]
      + fetches a URL and returns the response body
      - returns error on network failure or non-2xx status
      # networking
  std.text
    std.text.normalize_whitespace
      fn (input: string) -> string
      + collapses runs of whitespace and trims the result
      # text

browser_agent
  browser_agent.new_session
    fn (headless: bool) -> result[session_state, string]
    + launches a browser session
    - returns error when no browser backend is available
    # lifecycle
  browser_agent.close_session
    fn (state: session_state) -> void
    + terminates the browser session
    # lifecycle
  browser_agent.navigate
    fn (state: session_state, url: string) -> result[session_state, string]
    + navigates to the URL and waits for the page to load
    - returns error on navigation failure
    # navigation
    -> std.net.http_get
  browser_agent.snapshot
    fn (state: session_state) -> page_snapshot
    + returns a structured snapshot of visible interactive elements with stable ids
    ? ids are derived from the DOM path and remain stable across re-snapshots when the path is unchanged
    # observation
    -> std.text.normalize_whitespace
  browser_agent.describe_element
    fn (snap: page_snapshot, element_id: string) -> result[element_info, string]
    + returns role, text, and attributes for the element
    - returns error when the id is unknown
    # observation
  browser_agent.click
    fn (state: session_state, element_id: string) -> result[session_state, string]
    + clicks the element identified by the snapshot id
    - returns error when the element is not found or not clickable
    # actions
  browser_agent.type_text
    fn (state: session_state, element_id: string, text: string) -> result[session_state, string]
    + types into an input element identified by the snapshot id
    - returns error when the element does not accept text input
    # actions
  browser_agent.select_option
    fn (state: session_state, element_id: string, value: string) -> result[session_state, string]
    + selects a dropdown option by value
    - returns error when the option does not exist
    # actions
  browser_agent.wait_for
    fn (state: session_state, selector: string, timeout_ms: i32) -> result[session_state, string]
    + waits until an element matching the selector appears
    - returns error when the timeout elapses
    # synchronization
  browser_agent.extract_text
    fn (snap: page_snapshot, element_id: string) -> result[string, string]
    + returns the normalized visible text of the element
    - returns error when the id is unknown
    # extraction
    -> std.text.normalize_whitespace
  browser_agent.screenshot
    fn (state: session_state) -> result[bytes, string]
    + captures a PNG screenshot of the current viewport
    - returns error when the session is closed
    # observation
