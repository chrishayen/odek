# Requirement: "a library for encoding and decoding unix timestamps to and from JSON"

Two project functions. The JSON representation is the integer seconds since the unix epoch.

std: (all units exist)

epoch_json
  epoch_json.encode_timestamp
    fn (seconds: i64) -> string
    + returns the JSON literal for the integer seconds value
    + returns "0" for the zero epoch
    # serialization
  epoch_json.decode_timestamp
    fn (raw: string) -> result[i64, string]
    + parses a JSON integer as seconds since the epoch
    - returns error on non-integer input
    - returns error on a JSON null
    # deserialization
