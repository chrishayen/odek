# Requirement: "a financial statement exchange client and response parser"

Builds a request document, sends it over HTTP, and parses the SGML/XML response into typed records. No CLI — the library is importable.

std
  std.http
    std.http.post
      fn (url: string, headers: map[string, string], body: bytes) -> result[bytes, string]
      + returns response body on 2xx
      - returns error on network failure or non-2xx status
      # http_client
  std.sgml
    std.sgml.parse
      fn (raw: string) -> result[sgml_tree, string]
      + parses loose SGML used by OFX 1.x into a tree
      - returns error on unbalanced tags
      # parsing

ofx
  ofx.build_request
    fn (user: string, password: string, fid: string, org: string, statement_type: string) -> string
    + returns a well-formed OFX request document with signon and statement blocks
    # request_building
  ofx.send
    fn (server_url: string, request: string) -> result[string, string]
    + returns the raw OFX response body
    - returns error when the server rejects the credentials
    # transport
    -> std.http.post
  ofx.parse_statement
    fn (raw_response: string) -> result[ofx_statement, string]
    + returns an ofx_statement with account, balance, and transaction list
    - returns error when the response is missing the statement block
    # parsing
    -> std.sgml.parse
  ofx.transactions
    fn (statement: ofx_statement) -> list[ofx_transaction]
    + returns the list of transactions in server order
    # access
