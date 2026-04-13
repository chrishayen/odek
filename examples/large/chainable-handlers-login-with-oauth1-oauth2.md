# Requirement: "chainable HTTP handlers for login flows with OAuth1 and OAuth2 authentication providers"

Full OAuth login middleware. The project layer exposes flow builders; the std layer supplies crypto, encoding, and HTTP primitives so the same primitives serve both protocols.

std
  std.http
    std.http.get
      @ (url: string, headers: map[string, string]) -> result[http_response, string]
      + issues a GET request and returns status, headers, and body
      - returns error on network failure
      # http_client
    std.http.post_form
      @ (url: string, headers: map[string, string], form: map[string, string]) -> result[http_response, string]
      + issues a POST with url-encoded body
      - returns error on non-2xx response
      # http_client
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + encodes bytes to standard base64
      # encoding
    std.encoding.url_encode
      @ (value: string) -> string
      + percent-encodes a string for use in URL query or form body
      # encoding
    std.encoding.parse_query_string
      @ (raw: string) -> map[string, string]
      + parses url-encoded key/value pairs
      # encoding
  std.crypto
    std.crypto.hmac_sha1
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA1 of data under key
      # cryptography
    std.crypto.random_bytes
      @ (count: i32) -> bytes
      + returns cryptographically random bytes of the requested length
      # cryptography
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

oauth_login
  oauth_login.oauth2_authorize_url
    @ (client_id: string, redirect_uri: string, scope: string, state: string, auth_endpoint: string) -> string
    + returns a well-formed OAuth2 authorization URL with the standard query parameters
    # oauth2
    -> std.encoding.url_encode
  oauth_login.oauth2_exchange_code
    @ (code: string, client_id: string, client_secret: string, redirect_uri: string, token_endpoint: string) -> result[oauth_token, string]
    + exchanges an authorization code for an access token
    - returns error when the provider responds with an error object
    # oauth2
    -> std.http.post_form
    -> std.json.parse_object
  oauth_login.oauth2_fetch_user
    @ (token: oauth_token, userinfo_endpoint: string) -> result[map[string, string], string]
    + fetches the userinfo document using the access token
    - returns error on non-2xx response
    # oauth2
    -> std.http.get
    -> std.json.parse_object
  oauth_login.oauth1_request_token
    @ (consumer_key: string, consumer_secret: string, callback: string, request_endpoint: string) -> result[oauth1_token, string]
    + performs the OAuth1 request-token step and returns an unauthorized token
    - returns error on provider failure
    # oauth1
    -> std.crypto.hmac_sha1
    -> std.crypto.random_bytes
    -> std.time.now_seconds
    -> std.encoding.base64_encode
    -> std.http.post_form
    -> std.encoding.parse_query_string
  oauth_login.oauth1_authorize_url
    @ (token: oauth1_token, authorize_endpoint: string) -> string
    + returns the user-facing authorize URL for the request token
    # oauth1
    -> std.encoding.url_encode
  oauth_login.oauth1_access_token
    @ (token: oauth1_token, verifier: string, consumer_key: string, consumer_secret: string, access_endpoint: string) -> result[oauth1_token, string]
    + completes the OAuth1 flow and returns an authorized access token
    - returns error when the verifier is rejected
    # oauth1
    -> std.crypto.hmac_sha1
    -> std.http.post_form
    -> std.encoding.parse_query_string
  oauth_login.with_state_check
    @ (inner: login_handler) -> login_handler
    + wraps a handler that verifies the state parameter on callback to prevent CSRF
    - rejects callbacks whose state does not match the one issued at authorize time
    # middleware
  oauth_login.new_state_token
    @ () -> string
    + returns a random opaque state token suitable for use in the OAuth2 state parameter
    # state
    -> std.crypto.random_bytes
    -> std.encoding.base64_encode
