# Requirement: "a webhook receiver that validates and parses events from hosted git providers"

Verifies provider-specific HMAC signatures and decodes payloads into a unified event shape.

std
  std.crypto
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256
      + returns 32 bytes
      # cryptography
    std.crypto.constant_time_eq
      @ (a: bytes, b: bytes) -> bool
      + returns true when inputs are equal in constant time
      - returns false when lengths differ
      # cryptography
  std.encoding
    std.encoding.hex_decode
      @ (s: string) -> result[bytes, string]
      + decodes a lowercase hex string
      - returns error on non-hex characters
      # encoding
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a flat string map
      - returns error on invalid JSON
      # serialization

webhooks
  webhooks.verify_signature
    @ (provider: string, secret: string, body: bytes, header: string) -> result[void, string]
    + returns ok when the provider-specific signature header matches
    - returns error when signature format is invalid
    - returns error when the computed digest does not match
    # signature_verification
    -> std.encoding.hex_decode
    -> std.crypto.hmac_sha256
    -> std.crypto.constant_time_eq
  webhooks.parse_event
    @ (provider: string, event_name: string, body: bytes) -> result[event, string]
    + returns a normalized event with type, repo, and actor fields
    - returns error on unknown provider
    - returns error on malformed payload
    # parsing
    -> std.json.parse_object
  webhooks.receive
    @ (provider: string, secret: string, event_name: string, body: bytes, signature_header: string) -> result[event, string]
    + verifies the signature then parses the payload in one call
    - returns error when verification fails
    # pipeline
