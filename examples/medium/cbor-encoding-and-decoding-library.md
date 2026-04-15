# Requirement: "a cbor encoding and decoding library"

Encode and decode CBOR values. The project layer exposes a tagged value type and pipeline entry points.

std: (all units exist)

cbor
  cbor.int
    fn (v: i64) -> cbor_value
    + wraps an integer as a CBOR value
    # constructor
  cbor.text
    fn (v: string) -> cbor_value
    + wraps a string as a CBOR value
    # constructor
  cbor.bytes
    fn (v: bytes) -> cbor_value
    + wraps a byte string as a CBOR value
    # constructor
  cbor.array
    fn (items: list[cbor_value]) -> cbor_value
    + wraps a list as a CBOR value
    # constructor
  cbor.map
    fn (pairs: list[cbor_pair]) -> cbor_value
    + wraps key-value pairs as a CBOR map
    ? key ordering is preserved for deterministic output
    # constructor
  cbor.encode
    fn (value: cbor_value) -> bytes
    + returns the canonical CBOR byte sequence for the value
    # encoding
  cbor.decode
    fn (data: bytes) -> result[cbor_value, string]
    + parses CBOR bytes into a tagged value
    - returns error on truncated input
    - returns error on an unsupported major type
    # decoding
