# Requirement: "a UI testing framework driving a browser through a remote automation protocol"

Provides a test runner, an assertion library, and a typed browser control surface. Std holds HTTP and JSON; the project splits into driver, session, and runner.

std
  std.http
    std.http.request
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an HTTP request
      - returns error on network failure
      # http_client
  std.json
    std.json.encode
      @ (value: json_value) -> string
      + encodes a JSON value as text
      # serialization
    std.json.decode
      @ (raw: string) -> result[json_value, string]
      + parses text as a JSON value
      - returns error on malformed input
      # serialization
  std.time
    std.time.sleep_ms
      @ (ms: i32) -> void
      + blocks the caller for ms milliseconds
      # time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

ui_test
  ui_test.open_session
    @ (remote_url: string, browser: string) -> result[session, string]
    + starts a new browser session on the remote driver
    - returns error when the remote rejects the request
    # session
    -> std.http.request
    -> std.json.encode
    -> std.json.decode
  ui_test.close_session
    @ (sess: session) -> result[void, string]
    + shuts down the browser session
    # session
    -> std.http.request
  ui_test.visit
    @ (sess: session, url: string) -> result[void, string]
    + navigates the browser to url
    # navigation
    -> std.http.request
  ui_test.locate
    @ (sess: session, selector: string) -> result[element, string]
    + returns a handle to the first element matching the CSS selector
    - returns error when no element matches
    # lookup
    -> std.http.request
  ui_test.click
    @ (sess: session, elem: element) -> result[void, string]
    + clicks the element
    # interaction
    -> std.http.request
  ui_test.type_text
    @ (sess: session, elem: element, text: string) -> result[void, string]
    + types text into the element
    # interaction
    -> std.http.request
    -> std.json.encode
  ui_test.wait_for
    @ (sess: session, selector: string, timeout_ms: i32) -> result[element, string]
    + polls locate until a matching element is visible
    - returns error when the timeout elapses
    # synchronization
    -> std.time.sleep_ms
    -> std.time.now_millis
  ui_test.assert_text
    @ (sess: session, elem: element, expected: string) -> result[void, string]
    + returns ok when the element's text equals expected
    - returns error with a diff message when the text differs
    # assertion
    -> std.http.request
  ui_test.register_test
    @ (suite: test_suite, name: string, body: test_body) -> test_suite
    + adds a named test to the suite
    # registration
  ui_test.run_suite
    @ (suite: test_suite, remote_url: string) -> suite_report
    + runs every registered test and returns pass/fail counts with per-test details
    # runner
  ui_test.render_report
    @ (report: suite_report) -> string
    + returns a text report listing failed tests with their messages
    # reporting
