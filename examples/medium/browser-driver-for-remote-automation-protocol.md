# Requirement: "a library for driving a remote browser via its automation protocol"

Talks to a browser automation endpoint over HTTP+JSON. Creates sessions, finds elements, and performs input actions.

std
  std.http
    std.http.post_json
      fn (url: string, body: string) -> result[http_response, string]
      + performs an HTTP POST with a JSON body
      - returns error on non-2xx response
      # network
    std.http.get
      fn (url: string) -> result[http_response, string]
      + performs an HTTP GET
      # network
    std.http.delete
      fn (url: string) -> result[http_response, string]
      + performs an HTTP DELETE
      # network
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object as a flat string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

browser_driver
  browser_driver.new
    fn (endpoint: string) -> driver_state
    + stores the automation server base URL
    # construction
  browser_driver.new_session
    fn (state: driver_state, capabilities: map[string, string]) -> result[tuple[driver_state, session_id], string]
    + starts a browser session and returns its id
    - returns error when the server rejects the capabilities
    # session
    -> std.http.post_json
    -> std.json.encode_object
    -> std.json.parse_object
  browser_driver.close_session
    fn (state: driver_state, id: session_id) -> result[driver_state, string]
    + closes the browser session
    # session
    -> std.http.delete
  browser_driver.navigate
    fn (state: driver_state, id: session_id, url: string) -> result[void, string]
    + navigates the session to the URL
    - returns error on navigation failure
    # navigation
    -> std.http.post_json
  browser_driver.find_element
    fn (state: driver_state, id: session_id, strategy: string, selector: string) -> result[element_handle, string]
    + locates a single element by strategy ("css", "xpath", "id")
    - returns error when no element matches
    # query
    -> std.http.post_json
    -> std.json.parse_object
  browser_driver.click
    fn (state: driver_state, id: session_id, element: element_handle) -> result[void, string]
    + dispatches a click at the element center
    - returns error when the element is stale
    # input
    -> std.http.post_json
  browser_driver.send_keys
    fn (state: driver_state, id: session_id, element: element_handle, text: string) -> result[void, string]
    + types text into the element
    # input
    -> std.http.post_json
  browser_driver.get_text
    fn (state: driver_state, id: session_id, element: element_handle) -> result[string, string]
    + returns the rendered text content of the element
    # query
    -> std.http.get
    -> std.json.parse_object
