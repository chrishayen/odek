# Requirement: "a standard envelope for wrapping proto messages published to message brokers"

Wraps a payload with routing metadata so downstream consumers have a uniform envelope regardless of broker.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.uuid
    std.uuid.v4
      @ () -> string
      + returns a random uuid in 8-4-4-4-12 hex form
      # identifiers

proto_envelope
  proto_envelope.wrap
    @ (payload: bytes, message_type: string, source: string, headers: map[string, string]) -> envelope
    + returns an envelope with a fresh id, current timestamp, and the given fields
    + copies headers so later mutation of the input map does not leak in
    # packaging
    -> std.uuid.v4
    -> std.time.now_millis
  proto_envelope.encode
    @ (e: envelope) -> bytes
    + serializes the envelope and its embedded payload
    # serialization
  proto_envelope.decode
    @ (data: bytes) -> result[envelope, string]
    + returns the parsed envelope
    - returns error on a missing required field
    - returns error on a truncated payload
    # deserialization
