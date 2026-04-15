# Requirement: "a browser automation client speaking the WebDriver protocol"

Wraps the WebDriver JSON-over-HTTP protocol as a typed session API. Std holds the HTTP and JSON primitives; the project exposes session, element lookup, interaction, and assertions.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an HTTP request and returns the response
      - returns error on network failure
      # http_client
  std.json
    std.json.encode
      fn (value: json_value) -> string
      + encodes a generic JSON value as text
      # serialization
    std.json.decode
      fn (raw: string) -> result[json_value, string]
      + parses text into a generic JSON value
      - returns error on malformed JSON
      # serialization
    std.json.get_string
      fn (value: json_value, path: string) -> result[string, string]
      + returns the string at a dotted path inside a JSON value
      - returns error when the path is missing or not a string
      # serialization

webdriver
  webdriver.new_session
    fn (base_url: string, capabilities: map[string, string]) -> result[session, string]
    + opens a WebDriver session and returns its identifier
    - returns error when the server rejects the capabilities
    # session
    -> std.http.request
    -> std.json.encode
    -> std.json.decode
  webdriver.quit
    fn (sess: session) -> result[void, string]
    + terminates the WebDriver session
    # session
    -> std.http.request
  webdriver.navigate
    fn (sess: session, url: string) -> result[void, string]
    + instructs the driver to load url
    - returns error when navigation fails
    # navigation
    -> std.http.request
    -> std.json.encode
  webdriver.find_element
    fn (sess: session, strategy: string, selector: string) -> result[element, string]
    + returns an element handle for the first match
    - returns error when no element matches
    # lookup
    -> std.http.request
    -> std.json.get_string
  webdriver.find_elements
    fn (sess: session, strategy: string, selector: string) -> result[list[element], string]
    + returns element handles for all matches, possibly empty
    # lookup
    -> std.http.request
  webdriver.click
    fn (sess: session, elem: element) -> result[void, string]
    + clicks the element
    - returns error when the element is not interactable
    # interaction
    -> std.http.request
  webdriver.send_keys
    fn (sess: session, elem: element, text: string) -> result[void, string]
    + types text into the element
    # interaction
    -> std.http.request
    -> std.json.encode
  webdriver.get_text
    fn (sess: session, elem: element) -> result[string, string]
    + returns the visible text of the element
    # query
    -> std.http.request
    -> std.json.get_string
  webdriver.get_attribute
    fn (sess: session, elem: element, name: string) -> result[optional[string], string]
    + returns the attribute value, or none when the attribute is absent
    # query
    -> std.http.request
  webdriver.wait_until_visible
    fn (sess: session, strategy: string, selector: string, timeout_ms: i32) -> result[element, string]
    + polls find_element until the element is present and displayed
    - returns error when the timeout elapses
    # synchronization
  webdriver.screenshot
    fn (sess: session) -> result[bytes, string]
    + captures the current viewport as PNG bytes
    # capture
    -> std.http.request
