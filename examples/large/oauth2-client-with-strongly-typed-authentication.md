# Requirement: "a strongly-typed OAuth2 client library"

Covers the authorization code, client credentials, and refresh flows. The project layer wires std primitives for URL building, JSON, and HTTP.

std
  std.url
    std.url.encode_query
      @ (params: map[string, string]) -> string
      + encodes parameters as an application/x-www-form-urlencoded string
      # url
    std.url.build
      @ (base: string, params: map[string, string]) -> string
      + appends an encoded query string to base, using ? or & as appropriate
      # url
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + encodes bytes to standard base64 with padding
      # encoding
  std.random
    std.random.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # random
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.http
    std.http.post_form
      @ (url: string, headers: map[string, string], body: string) -> result[http_response, string]
      + performs an HTTP POST with form-encoded body and returns the response
      - returns error on transport failure
      # http

oauth2
  oauth2.new_client
    @ (client_id: string, client_secret: string, auth_url: string, token_url: string, redirect_url: string) -> client_config
    + creates a client configuration holding endpoint URLs and credentials
    # construction
  oauth2.authorization_url
    @ (client: client_config, scopes: list[string], state: string) -> string
    + returns the authorization endpoint URL with response_type, client_id, redirect_uri, scope, and state parameters
    # authorization_code
    -> std.url.build
  oauth2.new_state
    @ () -> string
    + returns a fresh opaque state token suitable for CSRF protection
    # authorization_code
    -> std.random.random_bytes
    -> std.encoding.base64_encode
  oauth2.exchange_code
    @ (client: client_config, code: string) -> result[token, string]
    + exchanges an authorization code for a token using the token endpoint
    - returns error when the server returns a non-success status
    - returns error when the response is missing access_token
    # authorization_code
    -> std.url.encode_query
    -> std.http.post_form
    -> std.json.parse_object
  oauth2.client_credentials
    @ (client: client_config, scopes: list[string]) -> result[token, string]
    + requests a token using grant_type=client_credentials
    - returns error when the server returns a non-success status
    # client_credentials
    -> std.url.encode_query
    -> std.http.post_form
    -> std.json.parse_object
  oauth2.refresh
    @ (client: client_config, refresh_token: string) -> result[token, string]
    + requests a new token using grant_type=refresh_token
    - returns error when the refresh token is rejected
    # refresh
    -> std.url.encode_query
    -> std.http.post_form
    -> std.json.parse_object
  oauth2.is_expired
    @ (t: token, now_seconds: i64) -> bool
    + returns true when the token's expiry is at or before now_seconds
    # expiry
  oauth2.basic_auth_header
    @ (client_id: string, client_secret: string) -> string
    + returns the value of an HTTP Basic Authorization header for the credentials
    # authorization_code
    -> std.encoding.base64_encode
