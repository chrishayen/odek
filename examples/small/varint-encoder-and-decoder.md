# Requirement: "a variable-length integer encoder and decoder"

Encodes unsigned integers using 7 bits per byte with a continuation bit; signed integers use zig-zag mapping before unsigned encoding.

std: (all units exist)

varint
  varint.encode_u64
    fn (value: u64) -> bytes
    + returns 1 byte for values below 128
    + returns up to 10 bytes for the maximum u64
    # encoding
  varint.decode_u64
    fn (data: bytes) -> result[tuple[u64, i32], string]
    + returns (value, bytes_consumed) on success
    - returns error when no terminating byte is found within 10 bytes
    - returns error when data is empty
    # decoding
  varint.encode_i64
    fn (value: i64) -> bytes
    + encodes via zig-zag so small negatives also pack compactly
    # encoding
  varint.decode_i64
    fn (data: bytes) -> result[tuple[i64, i32], string]
    + returns (value, bytes_consumed) on success
    - returns error on truncated input
    # decoding
