# Requirement: "a framework for declaring interactive dashboards as a linear script of widgets"

Widgets are appended to a page; inputs are read from a request-scoped state; rendering produces an HTML document.

std
  std.http
    std.http.parse_form
      @ (body: bytes) -> map[string, string]
      + parses urlencoded form bodies into a string-to-string map
      # http

dashboards
  dashboards.new_page
    @ (title: string) -> page_state
    + returns an empty page with the given title
    # construction
  dashboards.add_heading
    @ (page: page_state, text: string) -> page_state
    + appends a heading widget to the page
    # widgets
  dashboards.add_text
    @ (page: page_state, body: string) -> page_state
    + appends a paragraph widget with the given markdown-compatible body
    # widgets
  dashboards.add_text_input
    @ (page: page_state, name: string, label: string) -> page_state
    + appends a text input widget bound to the given name
    # widgets
  dashboards.add_slider
    @ (page: page_state, name: string, min: f64, max: f64) -> page_state
    + appends a numeric slider bound to the given name
    # widgets
  dashboards.add_chart
    @ (page: page_state, name: string, data: list[f64]) -> page_state
    + appends a line chart widget rendering the given values
    # widgets
  dashboards.read_input
    @ (inputs: map[string, string], name: string) -> optional[string]
    + returns the current value of the named input or none if unset
    # state
  dashboards.render
    @ (page: page_state, inputs: map[string, string]) -> string
    + returns a full HTML document for the page using the given input values
    # rendering
  dashboards.handle_submit
    @ (body: bytes) -> map[string, string]
    + parses a form submission into an inputs map for the next render
    # transport
    -> std.http.parse_form
