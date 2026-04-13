# Requirement: "a headless browser test scripting wrapper"

A thin DSL over a headless browser driver that makes writing navigation-and-assert test flows concise.

std
  std.time
    std.time.sleep_ms
      @ (ms: i32) -> void
      + blocks the caller for the given milliseconds
      # time

browser_test
  browser_test.new_session
    @ (start_url: string) -> result[session_state, string]
    + opens a headless browser pointed at start_url
    - returns error when the driver cannot be launched
    # session
  browser_test.visit
    @ (state: session_state, url: string) -> result[session_state, string]
    + navigates to url and waits for the document to be ready
    - returns error when the page fails to load within the timeout
    # navigation
  browser_test.fill
    @ (state: session_state, selector: string, value: string) -> result[session_state, string]
    + types value into the element matching selector
    - returns error when no element matches
    # interaction
  browser_test.click
    @ (state: session_state, selector: string) -> result[session_state, string]
    + clicks the element matching selector
    - returns error when no element matches
    # interaction
  browser_test.expect_text
    @ (state: session_state, selector: string, expected: string) -> result[void, string]
    + passes when the element text equals expected
    - returns error with actual text when it does not match
    # assertion
    -> std.time.sleep_ms
  browser_test.screenshot
    @ (state: session_state, path: string) -> result[void, string]
    + writes a PNG of the current viewport to path
    # capture
  browser_test.close
    @ (state: session_state) -> void
    + releases the browser and associated resources
    # teardown
