# Requirement: "a low- and high-level MessagePack encoder and decoder"

std: (all units exist)

msgpack
  msgpack.encode_nil
    @ () -> bytes
    + returns the single-byte nil encoding
    # low_level
  msgpack.encode_bool
    @ (b: bool) -> bytes
    + returns the single-byte true/false encoding
    # low_level
  msgpack.encode_int
    @ (n: i64) -> bytes
    + returns the shortest valid integer encoding (fixint, int8/16/32/64, uint variants)
    # low_level
  msgpack.encode_string
    @ (s: string) -> bytes
    + returns the shortest valid string encoding (fixstr, str8/16/32)
    # low_level
  msgpack.encode_array_header
    @ (n: i32) -> bytes
    + returns the shortest array-header encoding for n elements
    # low_level
  msgpack.encode_map_header
    @ (n: i32) -> bytes
    + returns the shortest map-header encoding for n entries
    # low_level
  msgpack.decode_value
    @ (data: bytes, offset: i32) -> result[tuple[msgpack_value, i32], string]
    + parses the value at offset and returns (value, next_offset)
    - returns error on truncated input or unknown type byte
    # low_level
  msgpack.encode_value
    @ (v: msgpack_value) -> bytes
    + recursively encodes a tagged value tree using the shortest forms
    # high_level
  msgpack.decode_all
    @ (data: bytes) -> result[msgpack_value, string]
    + decodes the entire buffer and returns the root value
    - returns error when trailing bytes remain after the root value
    # high_level
