# Requirement: "a common client library for document-oriented databases with a pluggable driver interface"

A client dispatches document CRUD and view queries to a registered driver. The project layer defines the uniform API and delegates operations through driver identifiers.

std
  std.http
    std.http.get
      @ (url: string, headers: map[string, string]) -> result[bytes, string]
      + sends a GET request and returns the body
      - returns error on non-2xx status
      # http
    std.http.put
      @ (url: string, headers: map[string, string], body: bytes) -> result[bytes, string]
      + sends a PUT request
      - returns error on non-2xx status
      # http
    std.http.delete
      @ (url: string, headers: map[string, string]) -> result[void, string]
      + sends a DELETE request
      - returns error on non-2xx status
      # http
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object
      - returns error on invalid input
      # serialization

doc_client
  doc_client.register_driver
    @ (name: string, driver_id: string) -> void
    + registers a driver under a name
    # driver_registration
  doc_client.new
    @ (driver_name: string, endpoint: string) -> result[client_state, string]
    + creates a client backed by a registered driver
    - returns error when the driver name is not registered
    # construction
  doc_client.create_database
    @ (state: client_state, name: string) -> result[void, string]
    + creates a database
    - returns error when the database already exists
    # database_management
  doc_client.delete_database
    @ (state: client_state, name: string) -> result[void, string]
    + deletes a database
    - returns error when the database does not exist
    # database_management
    -> std.http.delete
  doc_client.put_document
    @ (state: client_state, database: string, id: string, doc: map[string, string]) -> result[string, string]
    + stores a document and returns its new revision
    - returns error on revision conflict
    # documents
    -> std.json.encode_object
    -> std.http.put
  doc_client.get_document
    @ (state: client_state, database: string, id: string) -> result[map[string, string], string]
    + returns the document by id
    - returns error when the document does not exist
    # documents
    -> std.http.get
    -> std.json.parse_object
  doc_client.delete_document
    @ (state: client_state, database: string, id: string, revision: string) -> result[void, string]
    + deletes a document at the given revision
    - returns error on revision conflict
    # documents
    -> std.http.delete
  doc_client.query_view
    @ (state: client_state, database: string, design: string, view: string, key: string) -> result[list[map[string, string]], string]
    + runs a view and returns matching rows
    - returns error when the view does not exist
    # views
    -> std.http.get
    -> std.json.parse_object
  doc_client.all_docs
    @ (state: client_state, database: string) -> result[list[string], string]
    + returns every document id in the database
    # query
    -> std.http.get
    -> std.json.parse_object
  doc_client.changes
    @ (state: client_state, database: string, since: string) -> result[list[map[string, string]], string]
    + returns the changes feed since the given sequence
    # replication
    -> std.http.get
    -> std.json.parse_object
  doc_client.close
    @ (state: client_state) -> void
    + releases any driver-held resources
    # teardown
