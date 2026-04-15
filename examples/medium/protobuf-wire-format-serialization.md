# Requirement: "an implementation of a compact wire-format serialization protocol"

Schema-driven record encoding with varint-based field framing. Std owns the varint primitive.

std
  std.encoding
    std.encoding.varint_encode
      fn (n: u64) -> bytes
      + returns the base-128 varint encoding of n
      + returns a single zero byte for n == 0
      # encoding
    std.encoding.varint_decode
      fn (data: bytes, offset: i32) -> result[tuple[u64, i32], string]
      + returns (value, new_offset) after consuming a varint
      - returns error on truncated input
      - returns error when the varint exceeds 10 bytes
      # encoding
    std.encoding.zigzag_encode
      fn (n: i64) -> u64
      + returns the zigzag encoding of a signed integer
      # encoding
    std.encoding.zigzag_decode
      fn (n: u64) -> i64
      + returns the signed integer from a zigzag-encoded unsigned
      # encoding

protobuf
  protobuf.schema_new
    fn () -> protobuf_schema
    + returns an empty schema
    # construction
  protobuf.schema_add_field
    fn (s: protobuf_schema, field_number: i32, name: string, type_tag: string) -> result[protobuf_schema, string]
    + returns a new schema with the field added
    - returns error when field_number is already used
    - returns error when type_tag is unknown
    # schema
  protobuf.encode
    fn (s: protobuf_schema, record: map[string, bytes]) -> result[bytes, string]
    + returns the wire bytes with each field preceded by a tag varint (field_number << 3 | wire_type)
    - returns error when record contains a field not in the schema
    # encoding
    -> std.encoding.varint_encode
    -> std.encoding.zigzag_encode
  protobuf.decode
    fn (s: protobuf_schema, data: bytes) -> result[map[string, bytes], string]
    + returns a map of field name to value bytes
    - returns error when a field number in the data is not in the schema
    - returns error on truncated input
    # decoding
    -> std.encoding.varint_decode
    -> std.encoding.zigzag_decode
