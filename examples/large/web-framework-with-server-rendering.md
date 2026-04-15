# Requirement: "a server-rendered web application framework"

Routes map request paths to handlers that return a rendered page. The framework wires templates, data, and an HTTP-shaped request/response together. Network IO is the caller's concern.

std
  std.text
    std.text.render_template
      fn (template: string, vars: map[string, string]) -> result[string, string]
      + substitutes {{key}} occurrences from vars
      - returns error when a referenced key is missing
      # templating
    std.text.split_path_segments
      fn (path: string) -> list[string]
      + splits on '/' and drops empty segments
      # text
  std.http
    std.http.parse_query
      fn (query: string) -> map[string, string]
      + parses "a=1&b=2" into a map
      + returns an empty map for empty input
      # http
    std.http.status_text
      fn (code: i32) -> string
      + returns the canonical reason phrase for a status code
      # http

web_framework
  web_framework.new_app
    fn () -> app_state
    + returns an app with no routes
    # construction
  web_framework.add_route
    fn (state: app_state, method: string, pattern: string, template: string) -> app_state
    + registers a template-rendering route
    ? pattern segments of the form ":name" become path parameters
    # routing
  web_framework.match_route
    fn (state: app_state, method: string, path: string) -> optional[tuple[string, map[string, string]]]
    + returns (template, path_params) when a route matches
    - returns empty when no route matches
    # routing
    -> std.text.split_path_segments
  web_framework.handle_request
    fn (state: app_state, method: string, path: string, query: string, data: map[string, string]) -> tuple[i32, string]
    + returns (status, body) with body rendered from the matched template
    + merges path params, query params, and data into the template context
    - returns (404, "Not Found") when no route matches
    - returns (500, <error>) when the template references a missing key
    # request_handling
    -> std.http.parse_query
    -> std.text.render_template
    -> std.http.status_text
  web_framework.render_page
    fn (state: app_state, template: string, data: map[string, string]) -> result[string, string]
    + renders a template by name without going through a request
    - returns error when the template name is not registered
    # rendering
    -> std.text.render_template
