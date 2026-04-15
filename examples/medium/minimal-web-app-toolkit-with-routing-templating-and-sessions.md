# Requirement: "a minimal web application toolkit with routing, templating, and sessions"

Three small subsystems a basic web app needs. Each is independent.

std
  std.http
    std.http.parse_request
      fn (raw: string) -> result[http_request, string]
      + parses method, path, headers, and body
      - returns error on malformed input
      # http
  std.id
    std.id.generate
      fn () -> string
      + returns a new unique identifier
      # id_generation

web_app
  web_app.route
    fn (routes: list[route_entry], method: string, path: string) -> optional[handler_fn]
    + returns the first matching handler
    - returns none when no route matches
    # routing
  web_app.render_template
    fn (template: string, values: map[string, string]) -> result[string, string]
    + substitutes "{{name}}" placeholders with values from the map
    - returns error when a referenced key is missing
    # templating
  web_app.new_session
    fn (sessions: session_store) -> tuple[string, session_store]
    + returns a new session id and updated store
    # sessions
    -> std.id.generate
  web_app.get_session
    fn (sessions: session_store, id: string) -> optional[map[string, string]]
    + returns the session data when present
    - returns none for unknown ids
    # sessions
  web_app.set_session_value
    fn (sessions: session_store, id: string, key: string, value: string) -> result[session_store, string]
    + stores a value in the session
    - returns error for unknown session ids
    # sessions
