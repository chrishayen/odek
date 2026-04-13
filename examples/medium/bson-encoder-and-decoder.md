# Requirement: "a BSON encoder and decoder"

BSON documents are an ordered map of typed values. The project implements the type-tag framing; std provides the raw integer and string primitives.

std
  std.encoding
    std.encoding.encode_i32_le
      @ (value: i32) -> bytes
      + encodes an i32 as 4 little-endian bytes
      # encoding
    std.encoding.encode_i64_le
      @ (value: i64) -> bytes
      + encodes an i64 as 8 little-endian bytes
      # encoding
    std.encoding.encode_f64_le
      @ (value: f64) -> bytes
      + encodes an f64 as 8 little-endian bytes
      # encoding
    std.encoding.decode_i32_le
      @ (data: bytes, offset: i32) -> result[i32, string]
      + decodes 4 little-endian bytes at offset
      - returns error when offset+4 exceeds length
      # encoding
    std.encoding.decode_i64_le
      @ (data: bytes, offset: i32) -> result[i64, string]
      + decodes 8 little-endian bytes at offset
      - returns error when offset+8 exceeds length
      # encoding
    std.encoding.decode_f64_le
      @ (data: bytes, offset: i32) -> result[f64, string]
      + decodes 8 little-endian bytes at offset
      - returns error when offset+8 exceeds length
      # encoding

bson
  bson.encode_document
    @ (fields: list[tuple[string, bson_value]]) -> bytes
    + returns the full BSON document starting with total length and ending with null terminator
    # encoding
    -> std.encoding.encode_i32_le
  bson.encode_value
    @ (name: string, value: bson_value) -> bytes
    + returns the type byte, cstring name, and encoded payload for one element
    # encoding
    -> std.encoding.encode_i32_le
    -> std.encoding.encode_i64_le
    -> std.encoding.encode_f64_le
  bson.decode_document
    @ (data: bytes) -> result[list[tuple[string, bson_value]], string]
    + returns the ordered list of fields
    - returns error when the declared length does not match the buffer length
    - returns error when the terminator byte is missing
    # decoding
    -> std.encoding.decode_i32_le
  bson.decode_value
    @ (data: bytes, offset: i32) -> result[tuple[string, bson_value, i32], string]
    + returns the field name, value, and next offset for one element
    - returns error on unknown type tags
    # decoding
    -> std.encoding.decode_i32_le
    -> std.encoding.decode_i64_le
    -> std.encoding.decode_f64_le
  bson.read_cstring
    @ (data: bytes, offset: i32) -> result[tuple[string, i32], string]
    + returns the null-terminated string and the offset after the terminator
    - returns error when no null byte is found before the buffer end
    # decoding
