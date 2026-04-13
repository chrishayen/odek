# Requirement: "a library for managing and orchestrating a JWT-based authentication workflow with refresh tokens and revocation"

Handles login (issuing an access/refresh token pair), refresh (trading a refresh token for a new pair), revocation, and verification. A pluggable store keeps revocation and refresh state.

std
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
    std.encoding.base64url_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes base64url with or without padding
      - returns error on invalid alphabet
      # encoding
  std.crypto
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + returns HMAC-SHA256 of data under key
      # cryptography
  std.json
    std.json.encode_object
      @ (obj: map[string,string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.random
    std.random.bytes
      @ (n: i32) -> bytes
      + returns n cryptographically secure random bytes
      # random

jwtvault
  jwtvault.new
    @ (secret: bytes, access_ttl: i64, refresh_ttl: i64) -> jwtvault_state
    + creates a vault with the given secret and token lifetimes
    # construction
  jwtvault.issue
    @ (state: jwtvault_state, subject: string, claims: map[string,string]) -> result[token_pair, string]
    + returns a fresh access token and a fresh refresh token bound to the subject
    + records the refresh token as active in the store
    - returns error on empty subject
    # login
    -> std.time.now_seconds
    -> std.random.bytes
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
  jwtvault.sign_access
    @ (state: jwtvault_state, subject: string, claims: map[string,string], now: i64) -> string
    + builds and signs a short-lived access JWT
    # signing
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
  jwtvault.sign_refresh
    @ (state: jwtvault_state, subject: string, jti: string, now: i64) -> string
    + builds and signs a long-lived refresh JWT carrying a unique jti
    # signing
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
  jwtvault.verify_access
    @ (state: jwtvault_state, token: string) -> result[map[string,string], string]
    + returns the claim map when the signature is valid and the token is not expired
    - returns error on bad signature
    - returns error when exp is in the past
    # verification
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.json.parse_object
    -> std.time.now_seconds
  jwtvault.verify_refresh
    @ (state: jwtvault_state, token: string) -> result[refresh_claims, string]
    + returns the subject and jti when valid and the jti is still active in the store
    - returns error when the jti has been revoked
    # verification
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.json.parse_object
    -> std.time.now_seconds
  jwtvault.refresh
    @ (state: jwtvault_state, refresh_token: string) -> result[token_pair, string]
    + revokes the presented refresh token and issues a new token pair for the same subject
    - returns error when the refresh token is invalid, expired, or already revoked
    # rotation
  jwtvault.revoke
    @ (state: jwtvault_state, jti: string) -> result[void, string]
    + marks the given refresh jti as revoked in the store
    # revocation
  jwtvault.revoke_all_for_subject
    @ (state: jwtvault_state, subject: string) -> result[i32, string]
    + revokes every active refresh jti for the subject and returns the count revoked
    # revocation
