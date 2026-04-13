# Requirement: "a webhook receiver that validates, transforms, and persists payloads"

Receives signed payloads, verifies them, applies a user-defined transform, and stores them in a pluggable sink.

std
  std.crypto
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256 of data under key
      # cryptography
    std.crypto.constant_time_eq
      @ (a: bytes, b: bytes) -> bool
      + returns true when inputs are byte-equal in constant time
      # cryptography
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding
    std.encoding.hex_decode
      @ (s: string) -> result[bytes, string]
      + decodes a hex string
      - returns error on odd length or non-hex characters
      # encoding
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses any JSON value into a tagged union
      - returns error on invalid JSON
      # serialization
    std.json.encode_value
      @ (v: json_value) -> string
      + serializes a JSON value to a compact string
      # serialization
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

webhook
  webhook.new_config
    @ (secret: string, header_name: string) -> config
    + creates a receiver configuration
    # configuration
  webhook.verify
    @ (cfg: config, body: bytes, signature: string) -> result[void, string]
    + returns ok when signature matches the body
    - returns error on invalid signature format
    - returns error on digest mismatch
    # signature_verification
    -> std.encoding.hex_decode
    -> std.crypto.hmac_sha256
    -> std.crypto.constant_time_eq
  webhook.new_transform
    @ (field_map: map[string, string]) -> transform
    + creates a transform that renames and filters top-level fields
    # transformation
  webhook.apply_transform
    @ (t: transform, value: json_value) -> json_value
    + returns the transformed payload
    # transformation
  webhook.build_record
    @ (payload: json_value, received_at: i64) -> stored_record
    + wraps a payload with timestamp and a derived id
    # recording
    -> std.time.now_millis
  webhook.store
    @ (sink: record_sink, record: stored_record) -> result[void, string]
    + persists a record via the provided sink
    - returns error when the sink rejects the record
    # persistence
    -> std.json.encode_value
  webhook.handle
    @ (cfg: config, t: transform, sink: record_sink, body: bytes, signature: string) -> result[void, string]
    + verifies, parses, transforms, and stores a payload
    - returns error at any failing stage
    # pipeline
    -> std.json.parse_value
