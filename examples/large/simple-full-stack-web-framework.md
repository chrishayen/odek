# Requirement: "a full-stack web framework"

Routing, request handling, templating, and a small ORM-lite data layer integrated into a single application object.

std
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses method, path, headers, and body from raw bytes
      - returns error on malformed request
      # http
    std.http.encode_response
      @ (status: i32, headers: map[string, string], body: bytes) -> bytes
      + returns a wire-format HTTP response
      # http
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns file contents as a string
      - returns error when file does not exist
      # filesystem
  std.sql
    std.sql.execute
      @ (conn: db_conn, sql: string, params: list[string]) -> result[i64, string]
      + executes a statement and returns affected rows
      - returns error on sql error
      # sql
    std.sql.query
      @ (conn: db_conn, sql: string, params: list[string]) -> result[list[map[string, string]], string]
      + returns rows as string-to-string maps
      - returns error on sql error
      # sql

webframework
  webframework.new_app
    @ (db: db_conn) -> app_state
    + returns an empty application bound to the given database connection
    # construction
  webframework.route
    @ (state: app_state, method: string, path_pattern: string, handler: fn(http_request, app_state) -> http_response) -> app_state
    + registers a handler for method and path pattern
    + supports colon-prefixed path parameters
    # routing
  webframework.match_route
    @ (state: app_state, method: string, path: string) -> optional[route_match]
    + returns the matching route and extracted parameters
    - returns none when no route matches
    # routing
  webframework.handle_request
    @ (state: app_state, raw: bytes) -> result[bytes, string]
    + parses the request, dispatches to the handler, and returns the wire response
    - returns a 404 response when no route matches
    - returns a 500 response when the handler raises
    # dispatch
    -> std.http.parse_request
    -> std.http.encode_response
    -> webframework.match_route
  webframework.load_template
    @ (state: app_state, name: string) -> result[template, string]
    + reads and parses a template from the configured template directory
    - returns error when the template file is missing
    # templating
    -> std.fs.read_all
  webframework.render_template
    @ (tmpl: template, values: map[string, string]) -> result[string, string]
    + substitutes {{name}} placeholders with values
    - returns error on unbalanced braces
    # templating
  webframework.find_records
    @ (state: app_state, table: string, column: string, value: string) -> result[list[map[string, string]], string]
    + returns rows whose column equals value
    - returns error on sql failure
    # data
    -> std.sql.query
  webframework.insert_record
    @ (state: app_state, table: string, record: map[string, string]) -> result[i64, string]
    + inserts a row and returns affected rows
    - returns error on sql failure
    # data
    -> std.sql.execute
  webframework.json_response
    @ (status: i32, body: string) -> http_response
    + returns a response with content-type application/json
    # helpers
  webframework.html_response
    @ (status: i32, body: string) -> http_response
    + returns a response with content-type text/html
    # helpers
