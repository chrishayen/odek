# Requirement: "a client library for a cloud service provider that signs requests, routes them to the right service endpoint, and parses structured responses"

Service operations are described in a registry; the project layer exposes generic invoke plus a credential/signing pipeline.

std
  std.http
    std.http.send_request
      @ (method: string, url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
      + performs an HTTPS request and returns the response
      - returns error on network failure
      # http
  std.crypto
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256
      # cryptography
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + computes SHA-256
      # cryptography
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      - returns error on invalid JSON
      # parsing
    std.json.encode
      @ (value: json_value) -> string
      # serialization
  std.time
    std.time.now_iso8601
      @ () -> string
      + returns the current time formatted as ISO 8601 UTC
      # time
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      # encoding

cloud_client
  cloud_client.new
    @ (access_key: string, secret_key: string, region: string) -> client_state
    + constructs a client bound to a region and credentials
    # construction
  cloud_client.register_service
    @ (state: client_state, service: string, endpoint: string, operations: list[string]) -> client_state
    + records how to reach a service and which operations it supports
    # registration
  cloud_client.canonical_request
    @ (method: string, url: string, headers: map[string,string], body: bytes) -> string
    + builds the canonical request string for signing
    # signing
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  cloud_client.signing_key
    @ (secret_key: string, date: string, region: string, service: string) -> bytes
    + derives the per-request signing key
    # signing
    -> std.crypto.hmac_sha256
  cloud_client.sign_request
    @ (state: client_state, service: string, method: string, url: string, headers: map[string,string], body: bytes) -> map[string,string]
    + returns headers with an Authorization signature added
    # signing
    -> std.time.now_iso8601
    -> cloud_client.canonical_request
    -> cloud_client.signing_key
  cloud_client.invoke
    @ (state: client_state, service: string, operation: string, params: map[string,string]) -> result[json_value, string]
    + signs and sends the request, parses the JSON response
    - returns error when service or operation is not registered
    - returns error on HTTP non-2xx response
    # invocation
    -> cloud_client.sign_request
    -> std.http.send_request
    -> std.json.encode
    -> std.json.parse
  cloud_client.parse_error_response
    @ (body: string) -> optional[service_error]
    + extracts error code and message from a structured error body
    # errors
    -> std.json.parse
