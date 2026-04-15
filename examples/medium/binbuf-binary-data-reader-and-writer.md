# Requirement: "a library for reading and writing structured binary data in byte buffers"

Typed reads and writes into a byte buffer at an explicit offset, with little-endian primitives.

std: (all units exist)

binbuf
  binbuf.new
    fn (size: i32) -> bytes
    + returns a zero-filled buffer of the given size
    - returns empty bytes when size is zero or negative
    # construction
  binbuf.write_u32_le
    fn (buf: bytes, offset: i32, value: u32) -> result[bytes, string]
    + writes four little-endian bytes at the offset
    - returns error when offset+4 exceeds buffer length
    # writing
  binbuf.read_u32_le
    fn (buf: bytes, offset: i32) -> result[u32, string]
    + reads four little-endian bytes at the offset
    - returns error when offset+4 exceeds buffer length
    # reading
  binbuf.write_i64_le
    fn (buf: bytes, offset: i32, value: i64) -> result[bytes, string]
    + writes eight little-endian bytes at the offset
    - returns error when offset+8 exceeds buffer length
    # writing
  binbuf.read_i64_le
    fn (buf: bytes, offset: i32) -> result[i64, string]
    + reads eight little-endian bytes at the offset
    - returns error when offset+8 exceeds buffer length
    # reading
  binbuf.write_string
    fn (buf: bytes, offset: i32, value: string) -> result[bytes, string]
    + writes the UTF-8 bytes of value at the offset
    - returns error when the string does not fit
    # writing
  binbuf.read_string
    fn (buf: bytes, offset: i32, length: i32) -> result[string, string]
    + reads length bytes at the offset as a UTF-8 string
    - returns error when the range is out of bounds
    - returns error on invalid UTF-8
    # reading
