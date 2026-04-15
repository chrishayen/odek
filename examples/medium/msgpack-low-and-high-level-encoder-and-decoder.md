# Requirement: "a low- and high-level MessagePack encoder and decoder"

std: (all units exist)

msgpack
  msgpack.encode_nil
    fn () -> bytes
    + returns the single-byte nil encoding
    # low_level
  msgpack.encode_bool
    fn (b: bool) -> bytes
    + returns the single-byte true/false encoding
    # low_level
  msgpack.encode_int
    fn (n: i64) -> bytes
    + returns the shortest valid integer encoding (fixint, int8/16/32/64, uint variants)
    # low_level
  msgpack.encode_string
    fn (s: string) -> bytes
    + returns the shortest valid string encoding (fixstr, str8/16/32)
    # low_level
  msgpack.encode_array_header
    fn (n: i32) -> bytes
    + returns the shortest array-header encoding for n elements
    # low_level
  msgpack.encode_map_header
    fn (n: i32) -> bytes
    + returns the shortest map-header encoding for n entries
    # low_level
  msgpack.decode_value
    fn (data: bytes, offset: i32) -> result[tuple[msgpack_value, i32], string]
    + parses the value at offset and returns (value, next_offset)
    - returns error on truncated input or unknown type byte
    # low_level
  msgpack.encode_value
    fn (v: msgpack_value) -> bytes
    + recursively encodes a tagged value tree using the shortest forms
    # high_level
  msgpack.decode_all
    fn (data: bytes) -> result[msgpack_value, string]
    + decodes the entire buffer and returns the root value
    - returns error when trailing bytes remain after the root value
    # high_level
