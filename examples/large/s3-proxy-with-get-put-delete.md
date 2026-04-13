# Requirement: "an authenticating object-storage proxy supporting GET, PUT and DELETE with OIDC and basic auth"

The proxy is a thin dispatcher in front of an object-store client and a pluggable authenticator. Auth and storage primitives belong in std.

std
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses an HTTP/1.1 request into method, path, headers, and body
      - returns error on malformed request line
      # http
    std.http.encode_response
      @ (status: i32, headers: map[string,string], body: bytes) -> bytes
      + serializes a response
      # http
  std.base64
    std.base64.decode
      @ (encoded: string) -> result[bytes, string]
      + decodes standard base64
      - returns error on invalid characters
      # encoding
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + computes SHA-256
      # cryptography
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256
      # cryptography
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses a JSON object to a string map
      - returns error on invalid JSON
      # serialization
  std.object_store
    std.object_store.get_object
      @ (bucket: string, key: string) -> result[bytes, string]
      + fetches an object from the backing store
      - returns error when the object does not exist
      # storage
    std.object_store.put_object
      @ (bucket: string, key: string, body: bytes) -> result[void, string]
      + writes an object to the backing store
      # storage
    std.object_store.delete_object
      @ (bucket: string, key: string) -> result[void, string]
      + removes an object
      - returns error when the object does not exist
      # storage

obj_proxy
  obj_proxy.new
    @ (bucket: string, basic_users: map[string,string], oidc_issuer: string, oidc_jwks: bytes) -> obj_proxy_state
    + constructs a proxy bound to one bucket with both basic and OIDC auth enabled
    ? passing an empty users map disables basic auth; passing an empty issuer disables OIDC
    # construction
  obj_proxy.authenticate_basic
    @ (state: obj_proxy_state, header: string) -> result[string, string]
    + returns the authenticated user when credentials match
    - returns error on missing, malformed, or unknown credentials
    # auth
    -> std.base64.decode
  obj_proxy.authenticate_bearer
    @ (state: obj_proxy_state, header: string) -> result[string, string]
    + returns the subject claim when the JWT signature verifies against the configured JWKS
    - returns error on expired or malformed tokens
    # auth
    -> std.crypto.hmac_sha256
    -> std.json.parse_object
  obj_proxy.handle_get
    @ (state: obj_proxy_state, key: string, auth_header: string) -> bytes
    + returns a 200 response with the object body after auth succeeds
    - returns 401 when neither auth method accepts the header
    - returns 404 when the object is missing
    # request_handling
    -> std.object_store.get_object
    -> std.http.encode_response
  obj_proxy.handle_put
    @ (state: obj_proxy_state, key: string, body: bytes, auth_header: string) -> bytes
    + returns 200 after writing the object
    - returns 401 on failed auth
    # request_handling
    -> std.object_store.put_object
    -> std.http.encode_response
  obj_proxy.handle_delete
    @ (state: obj_proxy_state, key: string, auth_header: string) -> bytes
    + returns 200 after deleting the object
    - returns 401 on failed auth
    - returns 404 when the object does not exist
    # request_handling
    -> std.object_store.delete_object
    -> std.http.encode_response
  obj_proxy.dispatch
    @ (state: obj_proxy_state, raw_request: bytes) -> bytes
    + routes a parsed request to the method-specific handler and returns the response bytes
    - returns 405 for methods other than GET, PUT, DELETE
    - returns 400 for unparseable requests
    # routing
    -> std.http.parse_request
    -> std.http.encode_response
