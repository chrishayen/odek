# Requirement: "a SQL driver that speaks the Avatica wire protocol to a remote SQL gateway"

A client library that opens connections, executes statements, and iterates result sets over the Avatica protocol. The wire format is protobuf over HTTP.

std
  std.http
    std.http.post
      fn (url: string, headers: map[string,string], body: bytes) -> result[bytes, string]
      + returns response body on 2xx
      - returns error on non-2xx or transport failure
      # http
  std.encoding
    std.encoding.protobuf_encode
      fn (message: map[string,bytes]) -> bytes
      + encodes a tag->wire-value map into protobuf bytes
      # serialization
    std.encoding.protobuf_decode
      fn (raw: bytes) -> result[map[string,bytes], string]
      + decodes protobuf bytes into a tag->wire-value map
      - returns error on truncated or malformed input
      # serialization
  std.uuid
    std.uuid.new_v4
      fn () -> string
      + returns a random 128-bit identifier as a hyphenated string
      # identity

avatica
  avatica.open
    fn (endpoint: string) -> result[avatica_conn, string]
    + returns a connection bound to the remote endpoint and a fresh connection id
    - returns error when the endpoint is not a valid url
    # connection
    -> std.uuid.new_v4
  avatica.close
    fn (conn: avatica_conn) -> result[void, string]
    + sends CloseConnectionRequest and releases local state
    # connection
    -> std.encoding.protobuf_encode
    -> std.http.post
  avatica.create_statement
    fn (conn: avatica_conn) -> result[avatica_stmt, string]
    + returns a statement handle valid for the lifetime of the connection
    - returns error when the connection is closed
    # statement
    -> std.encoding.protobuf_encode
    -> std.http.post
  avatica.execute
    fn (stmt: avatica_stmt, sql: string) -> result[avatica_result, string]
    + runs the query and returns a result handle holding the first frame of rows
    - returns error when the server returns an error response
    # execution
    -> std.encoding.protobuf_encode
    -> std.encoding.protobuf_decode
    -> std.http.post
  avatica.fetch_next
    fn (result: avatica_result) -> result[optional[list[list[string]]], string]
    + returns the next frame of rows, or none when the cursor is exhausted
    - returns error when the server reports the statement has been invalidated
    # paging
    -> std.encoding.protobuf_encode
    -> std.encoding.protobuf_decode
    -> std.http.post
  avatica.columns
    fn (result: avatica_result) -> list[string]
    + returns the column names of the result in order
    # metadata
  avatica.prepare
    fn (conn: avatica_conn, sql: string) -> result[avatica_prepared, string]
    + returns a prepared statement handle with parameter metadata
    - returns error on syntactically invalid sql
    # prepared_statement
    -> std.encoding.protobuf_encode
    -> std.encoding.protobuf_decode
    -> std.http.post
  avatica.execute_prepared
    fn (prepared: avatica_prepared, params: list[string]) -> result[avatica_result, string]
    + binds parameters positionally and executes
    - returns error when the parameter count does not match
    # execution
    -> std.encoding.protobuf_encode
    -> std.encoding.protobuf_decode
    -> std.http.post
  avatica.close_statement
    fn (stmt: avatica_stmt) -> result[void, string]
    + releases the statement on the server
    # statement
    -> std.encoding.protobuf_encode
    -> std.http.post
