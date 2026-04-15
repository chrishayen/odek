# Requirement: "an auto-generated admin panel with CRUD operations over arbitrary resources"

Resources are described by schemas; the library derives list, create, read, update, and delete endpoints and renders simple HTML forms.

std
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses JSON into a generic value tree
      - returns error on malformed input
      # serialization
    std.json.encode
      fn (value: json_value) -> string
      + serializes a generic value tree
      # serialization
  std.html
    std.html.escape
      fn (text: string) -> string
      + escapes characters unsafe in HTML text nodes
      # html
    std.html.render_form
      fn (fields: list[form_field]) -> string
      + produces a form element with inputs per field
      # html
  std.http
    std.http.parse_request
      fn (raw: bytes) -> result[http_request, string]
      + parses an HTTP/1.1 request
      - returns error on malformed input
      # parsing
    std.http.encode_response
      fn (status: i32, headers: map[string,string], body: bytes) -> bytes
      + serializes an HTTP/1.1 response
      # serialization

adminbro
  adminbro.resource_new
    fn (name: string, schema: resource_schema) -> resource_def
    + describes a resource and its fields
    # construction
  adminbro.panel_new
    fn () -> panel_state
    + creates an empty admin panel
    # construction
  adminbro.register
    fn (panel: panel_state, resource: resource_def, store: resource_store) -> panel_state
    + attaches a resource and its persistent store to the panel
    # registration
  adminbro.list_action
    fn (panel: panel_state, resource_name: string, query: map[string,string]) -> result[list[record], string]
    + returns the matching records for the resource
    # actions
  adminbro.get_action
    fn (panel: panel_state, resource_name: string, id: string) -> result[record, string]
    + fetches a single record by id
    - returns error when not found
    # actions
  adminbro.create_action
    fn (panel: panel_state, resource_name: string, values: map[string,string]) -> result[record, string]
    + validates and persists a new record
    - returns error when required fields are missing
    # actions
  adminbro.update_action
    fn (panel: panel_state, resource_name: string, id: string, values: map[string,string]) -> result[record, string]
    + merges the values into the existing record
    - returns error when the id is unknown
    # actions
  adminbro.delete_action
    fn (panel: panel_state, resource_name: string, id: string) -> result[void, string]
    + removes the record from the store
    # actions
  adminbro.render_list_page
    fn (panel: panel_state, resource_name: string, records: list[record]) -> string
    + renders an HTML table for the records
    # rendering
    -> std.html.escape
  adminbro.render_edit_page
    fn (panel: panel_state, resource_name: string, record: record) -> string
    + renders a form prefilled with the record's values
    # rendering
    -> std.html.render_form
  adminbro.handle
    fn (panel: panel_state, request_raw: bytes) -> result[bytes, string]
    + routes the HTTP request to the matching action and returns a response
    - returns error on malformed request
    # dispatch
    -> std.http.parse_request
    -> std.http.encode_response
    -> std.json.parse
    -> std.json.encode
