# Requirement: "a library for automating interaction with websites"

Models a headless browser session: follow links, fill forms, submit. HTTP and HTML parsing are std primitives.

std
  std.http
    std.http.get
      fn (url: string, headers: map[string, string]) -> result[string, string]
      + returns response body for 2xx status
      - returns error on non-2xx or network failure
      # network
    std.http.post_form
      fn (url: string, headers: map[string, string], form: map[string, string]) -> result[string, string]
      + posts a urlencoded form body and returns the response body
      - returns error on non-2xx
      # network
  std.html
    std.html.parse
      fn (raw: string) -> result[html_doc, string]
      + parses HTML into a navigable document tree
      - returns error on unrecoverable malformed input
      # parsing
    std.html.find_all
      fn (doc: html_doc, selector: string) -> list[html_node]
      + returns nodes matching a css selector
      # querying

browser
  browser.new
    fn () -> browser_state
    + creates a session with an empty cookie jar
    # construction
  browser.open
    fn (state: browser_state, url: string) -> result[browser_state, string]
    + fetches the url, stores the response as the current page, and returns updated state
    - returns error on non-2xx response
    # navigation
    -> std.http.get
    -> std.html.parse
  browser.follow_link
    fn (state: browser_state, link_text: string) -> result[browser_state, string]
    + finds the first anchor whose visible text matches and opens it
    - returns error when no matching link exists
    # navigation
    -> std.html.find_all
  browser.select_form
    fn (state: browser_state, name_or_id: string) -> result[browser_state, string]
    + marks a form as the active form for subsequent field operations
    - returns error when no form matches
    # forms
    -> std.html.find_all
  browser.set_field
    fn (state: browser_state, field_name: string, value: string) -> result[browser_state, string]
    + sets a value on the active form
    - returns error when the field does not exist on the active form
    # forms
  browser.submit
    fn (state: browser_state) -> result[browser_state, string]
    + submits the active form using its method and action, updating the current page
    - returns error when no form is active
    # forms
    -> std.http.post_form
  browser.current_html
    fn (state: browser_state) -> string
    + returns the raw HTML of the current page
    # introspection
