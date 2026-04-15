# Requirement: "a browser testing and automation library"

Drives a headless browser through opaque handles. The browser engine itself is a std primitive.

std
  std.browser
    std.browser.launch
      fn (headless: bool) -> result[browser_handle, string]
      + launches a browser process
      - returns error when no browser engine is available
      # browser
    std.browser.close
      fn (handle: browser_handle) -> void
      + terminates the browser process
      # browser
    std.browser.new_page
      fn (handle: browser_handle) -> result[page_handle, string]
      + opens a new tab
      # browser
    std.browser.goto
      fn (page: page_handle, url: string) -> result[page_handle, string]
      + navigates the page to the given URL
      - returns error on network failure
      # browser
    std.browser.query
      fn (page: page_handle, selector: string) -> optional[element_handle]
      + returns the first element matching the selector
      # browser
    std.browser.query_all
      fn (page: page_handle, selector: string) -> list[element_handle]
      + returns every element matching the selector
      # browser
    std.browser.click
      fn (element: element_handle) -> result[void, string]
      + dispatches a click on the element
      # browser
    std.browser.type_text
      fn (element: element_handle, text: string) -> result[void, string]
      + types the text into a focusable element
      # browser
    std.browser.inner_text
      fn (element: element_handle) -> string
      + returns the rendered text content
      # browser
    std.browser.screenshot
      fn (page: page_handle) -> result[bytes, string]
      + captures the current viewport as PNG bytes
      # browser

web
  web.open_browser
    fn (headless: bool) -> result[browser_state, string]
    + launches the browser and returns a driver state
    # lifecycle
    -> std.browser.launch
  web.close_browser
    fn (state: browser_state) -> void
    + shuts down the browser and releases resources
    # lifecycle
    -> std.browser.close
  web.new_page
    fn (state: browser_state) -> result[page_state, string]
    + opens a fresh page
    # navigation
    -> std.browser.new_page
  web.visit
    fn (page: page_state, url: string) -> result[page_state, string]
    + navigates to a URL and waits for the page to settle
    - returns error when navigation fails
    # navigation
    -> std.browser.goto
  web.click
    fn (page: page_state, selector: string) -> result[page_state, string]
    + clicks the first element matching the selector
    - returns error when no element matches
    # interaction
    -> std.browser.query
    -> std.browser.click
  web.fill
    fn (page: page_state, selector: string, text: string) -> result[page_state, string]
    + types text into the first matching element
    - returns error when no element matches
    # interaction
    -> std.browser.query
    -> std.browser.type_text
  web.text_content
    fn (page: page_state, selector: string) -> optional[string]
    + returns the rendered text of the first matching element
    # assertion
    -> std.browser.query
    -> std.browser.inner_text
  web.wait_for_selector
    fn (page: page_state, selector: string, timeout_ms: i32) -> result[page_state, string]
    + polls until a matching element appears or the timeout elapses
    - returns error when the timeout elapses
    # synchronization
    -> std.browser.query
  web.count
    fn (page: page_state, selector: string) -> i32
    + returns the number of elements matching the selector
    # assertion
    -> std.browser.query_all
  web.screenshot
    fn (page: page_state, path: string) -> result[void, string]
    + captures the viewport and writes it to path
    - returns error on capture or write failure
    # capture
    -> std.browser.screenshot
