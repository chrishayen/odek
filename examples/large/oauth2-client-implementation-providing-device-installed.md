# Requirement: "an OAuth2 client supporting device, installed-app, and service-account flows"

Three flows share the same token exchange and storage primitives. Real work lives in std crypto/http/json primitives.

std
  std.http
    std.http.post_form
      @ (url: string, form: map[string, string]) -> result[bytes, string]
      + posts url-encoded form and returns response body
      - returns error on non-2xx status
      # http
    std.http.get
      @ (url: string, headers: map[string, string]) -> result[bytes, string]
      + issues a GET with the given headers and returns the body
      - returns error on non-2xx status
      # http
  std.json
    std.json.parse_object
      @ (raw: bytes) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.crypto
    std.crypto.rsa_sign_sha256
      @ (private_key_pem: string, data: bytes) -> result[bytes, string]
      + signs data with an RSA private key using SHA-256
      - returns error when the key is malformed
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.url
    std.url.encode_query
      @ (params: map[string, string]) -> string
      + returns an rfc3986-encoded query string
      # encoding

oauth2
  oauth2.device_start
    @ (client_id: string, scopes: list[string], device_endpoint: string) -> result[map[string, string], string]
    + returns device_code, user_code, verification_uri, and interval
    - returns error when the endpoint rejects the request
    # device_flow
    -> std.http.post_form
    -> std.json.parse_object
  oauth2.device_poll
    @ (client_id: string, device_code: string, token_endpoint: string) -> result[oauth2_token, string]
    + returns an access token once the user has approved
    - returns error "authorization_pending" while waiting
    - returns error "access_denied" when the user declines
    # device_flow
    -> std.http.post_form
    -> std.json.parse_object
  oauth2.installed_authorize_url
    @ (client_id: string, redirect_uri: string, scopes: list[string], auth_endpoint: string) -> string
    + returns an authorization URL for a desktop/installed app with response_type=code
    # installed_flow
    -> std.url.encode_query
  oauth2.installed_exchange_code
    @ (client_id: string, client_secret: string, code: string, redirect_uri: string, token_endpoint: string) -> result[oauth2_token, string]
    + exchanges an authorization code for an access token
    - returns error when the code is invalid or expired
    # installed_flow
    -> std.http.post_form
    -> std.json.parse_object
  oauth2.service_account_sign_jwt
    @ (iss: string, scope: string, aud: string, private_key_pem: string) -> result[string, string]
    + builds and signs an RSA-SHA256 assertion JWT for service-account flow
    - returns error when the private key is malformed
    # service_account_flow
    -> std.encoding.base64url_encode
    -> std.crypto.rsa_sign_sha256
    -> std.time.now_seconds
  oauth2.service_account_exchange
    @ (assertion: string, token_endpoint: string) -> result[oauth2_token, string]
    + exchanges a signed assertion for an access token
    - returns error when the assertion is rejected
    # service_account_flow
    -> std.http.post_form
    -> std.json.parse_object
  oauth2.refresh
    @ (client_id: string, client_secret: string, refresh_token: string, token_endpoint: string) -> result[oauth2_token, string]
    + returns a fresh access token using a refresh token
    - returns error when the refresh token is revoked
    # refresh
    -> std.http.post_form
    -> std.json.parse_object
  oauth2.is_expired
    @ (token: oauth2_token) -> bool
    + returns true when the token's expiry is in the past
    # expiry
    -> std.time.now_seconds
