# Requirement: "a client driver for a document database with a fluent query language"

Connects to a server, authenticates, builds queries as composable term trees, and executes them over a pluggable transport.

std
  std.net
    std.net.dial_tcp
      @ (host: string, port: i32) -> result[conn_handle, string]
      + opens a TCP connection to host:port
      - returns error when the connection cannot be established
      # networking
    std.net.write_all
      @ (conn: conn_handle, data: bytes) -> result[void, string]
      + writes all bytes to the connection
      # networking
    std.net.read_exact
      @ (conn: conn_handle, n: i32) -> result[bytes, string]
      + reads exactly n bytes
      - returns error on short read
      # networking
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a map
      - returns error on invalid JSON
      # serialization

doc_db_driver
  doc_db_driver.connect
    @ (host: string, port: i32, auth_key: string) -> result[session_state, string]
    + opens a connection and performs the handshake with the auth key
    - returns error on handshake rejection
    # session
    -> std.net.dial_tcp
    -> std.net.write_all
    -> std.net.read_exact
  doc_db_driver.table
    @ (db_name: string, table_name: string) -> query_term
    + builds a table-reference term
    # query_builder
  doc_db_driver.filter
    @ (source: query_term, field: string, value: string) -> query_term
    + wraps source in a filter term matching field == value
    # query_builder
  doc_db_driver.limit
    @ (source: query_term, n: i32) -> query_term
    + wraps source in a limit term
    # query_builder
  doc_db_driver.run
    @ (session: session_state, term: query_term) -> result[list[string], string]
    + serializes the term tree, sends it, and returns JSON documents as strings
    - returns error on server-side failures
    # execution
    -> std.json.encode_object
    -> std.json.parse_object
    -> std.net.write_all
    -> std.net.read_exact
  doc_db_driver.close
    @ (session: session_state) -> result[void, string]
    + closes the underlying connection
    # session
