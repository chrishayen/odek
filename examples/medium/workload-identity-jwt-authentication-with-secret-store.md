# Requirement: "a library for workload-identity JWT authentication against a secret store"

A workload exchanges its identity JWT for a short-lived secret-store token, then fetches secrets with that token. Transport and JWT verification live in std.

std
  std.http
    std.http.post_json
      @ (url: string, body: string, headers: map[string, string]) -> result[tuple[i32, string], string]
      + sends a POST with JSON body and returns (status, response body)
      - returns error on network failure
      # http
    std.http.get
      @ (url: string, headers: map[string, string]) -> result[tuple[i32, string], string]
      + sends a GET and returns (status, response body)
      - returns error on network failure
      # http
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.jwt
    std.jwt.decode_unverified
      @ (token: string) -> result[map[string, string], string]
      + extracts claims without verifying the signature
      - returns error when the token does not have three segments
      # tokens

workload_secrets
  workload_secrets.extract_audience
    @ (svid_jwt: string) -> result[string, string]
    + reads the "aud" claim from an identity JWT
    - returns error when the claim is missing
    # claims
    -> std.jwt.decode_unverified
  workload_secrets.login
    @ (secret_store_url: string, role: string, svid_jwt: string) -> result[string, string]
    + posts the identity JWT to the login endpoint and returns the issued session token
    - returns error when the store rejects the token (non-2xx status)
    # authentication
    -> std.http.post_json
    -> std.json.encode_object
    -> std.json.parse_object
  workload_secrets.read_secret
    @ (secret_store_url: string, session_token: string, path: string) -> result[map[string, string], string]
    + fetches a secret at the given path and returns its key-value data
    - returns error when the session token is rejected
    - returns error when no secret exists at the path
    # retrieval
    -> std.http.get
    -> std.json.parse_object
