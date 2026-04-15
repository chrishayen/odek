# Requirement: "a headless CMS with an embedded web framework"

A full backend: content schema registry, content CRUD over HTTP, and a minimal routing layer. Real work sits in std (http server, json, sql).

std
  std.http
    std.http.serve
      fn (addr: string, handler: http_handler) -> result[void, string]
      + starts an HTTP server on the given address and dispatches to the handler
      - returns error when the address is already in use
      # http
    std.http.read_body
      fn (req: http_request) -> result[bytes, string]
      + reads the full request body
      # http
    std.http.write_response
      fn (resp: http_response, status: i32, body: bytes, content_type: string) -> void
      + writes status, content-type, and body to the response
      # http
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses arbitrary JSON
      - returns error on malformed input
      # serialization
    std.json.encode
      fn (value: json_value) -> string
      + encodes a json value to string
      # serialization
  std.sql
    std.sql.open
      fn (dsn: string) -> result[db_handle, string]
      + opens a database connection pool
      - returns error on unreachable database
      # database
    std.sql.exec
      fn (db: db_handle, query: string, params: list[json_value]) -> result[i64, string]
      + executes a statement and returns rows-affected
      # database
    std.sql.query
      fn (db: db_handle, query: string, params: list[json_value]) -> result[list[map[string, json_value]], string]
      + runs a query and returns rows as maps
      # database
  std.uuid
    std.uuid.v4
      fn () -> string
      + returns a new random UUID v4
      # identifier

cms
  cms.new
    fn (dsn: string) -> result[cms_state, string]
    + creates a cms bound to a database
    - returns error when the database cannot be opened
    # construction
    -> std.sql.open
  cms.register_schema
    fn (state: cms_state, type_name: string, fields: map[string, string]) -> result[cms_state, string]
    + registers a content type with named typed fields
    - returns error when a field uses an unknown type keyword
    # schema
  cms.create_entry
    fn (state: cms_state, type_name: string, values: map[string, json_value]) -> result[string, string]
    + validates values against the schema and inserts an entry, returning its id
    - returns error when a required field is missing or has the wrong type
    # content
    -> std.uuid.v4
    -> std.sql.exec
  cms.get_entry
    fn (state: cms_state, type_name: string, id: string) -> result[map[string, json_value], string]
    + returns the entry values
    - returns error when the id is not found
    # content
    -> std.sql.query
  cms.update_entry
    fn (state: cms_state, type_name: string, id: string, values: map[string, json_value]) -> result[void, string]
    + applies partial updates to an entry
    - returns error when the id is not found
    # content
    -> std.sql.exec
  cms.delete_entry
    fn (state: cms_state, type_name: string, id: string) -> result[void, string]
    + removes an entry
    - returns error when the id is not found
    # content
    -> std.sql.exec
  cms.list_entries
    fn (state: cms_state, type_name: string, limit: i32, offset: i32) -> result[list[map[string, json_value]], string]
    + returns a page of entries for a type
    # content
    -> std.sql.query
  cms.route_request
    fn (state: cms_state, req: http_request) -> http_response
    + dispatches REST routes /api/{type}, /api/{type}/{id} to the matching operation
    + returns 404 for unknown types
    # routing
    -> std.http.read_body
    -> std.json.parse
    -> std.json.encode
    -> std.http.write_response
  cms.serve
    fn (state: cms_state, addr: string) -> result[void, string]
    + starts the HTTP server bound to the cms router
    # server
    -> std.http.serve
