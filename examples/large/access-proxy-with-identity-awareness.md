# Requirement: "an identity-aware access proxy"

Forwards HTTP requests to upstreams only after authenticating the user and checking per-route authorization policies.

std
  std.crypto
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + returns a 32-byte MAC
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
    std.encoding.base64url_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes base64url with or without padding
      - returns error on invalid characters
      # encoding
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string,string]) -> string
      + encodes a string map as JSON
      # serialization
  std.net
    std.net.http_forward
      @ (upstream: string, req: http_request) -> result[http_response, string]
      + forwards req to upstream and returns the response
      - returns error when the upstream is unreachable
      # networking
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

access_proxy
  access_proxy.load_config
    @ (raw: string) -> result[config_state, string]
    + parses a proxy configuration document with routes, upstreams, and policies
    - returns error when a policy references an unknown claim
    # configuration
    -> std.json.parse_object
  access_proxy.issue_session
    @ (secret: bytes, subject: string, claims: map[string,string], ttl_seconds: i64) -> string
    + returns a signed session token binding subject and claims
    # sessions
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
    -> std.time.now_seconds
  access_proxy.verify_session
    @ (secret: bytes, token: string) -> result[session_claims, string]
    + returns claims when the token signature and expiry are valid
    - returns error when the signature does not match
    - returns error when the token is expired
    # sessions
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.json.parse_object
    -> std.time.now_seconds
  access_proxy.match_route
    @ (config: config_state, host: string, path: string) -> optional[route]
    + returns the first route whose host and path prefix match
    # routing
  access_proxy.evaluate_policy
    @ (route: route, claims: session_claims) -> result[void, string]
    + passes when every required claim predicate holds for the session
    - returns error naming the first failing predicate
    # authorization
  access_proxy.handle_request
    @ (state: proxy_state, req: http_request) -> http_response
    + redirects unauthenticated requests to the login flow
    + forwards authorized requests to the matched upstream
    - returns 403 when policy evaluation fails
    # dispatch
    -> std.net.http_forward
  access_proxy.start_login
    @ (state: proxy_state, return_url: string) -> http_response
    + returns a redirect to the identity provider with a state parameter
    # authn
  access_proxy.complete_login
    @ (state: proxy_state, code: string, state_param: string) -> result[http_response, string]
    + exchanges code for identity and sets the session cookie
    - returns error when state_param does not match
    # authn
  access_proxy.logout
    @ (state: proxy_state) -> http_response
    + returns a response that clears the session cookie
    # session_teardown
