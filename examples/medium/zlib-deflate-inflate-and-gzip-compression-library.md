# Requirement: "a zlib-compatible compression library supporting deflate, inflate, and gzip"

Implements the RFC 1950/1951/1952 stack over an in-memory buffer.

std: (all units exist)

zlib
  zlib.deflate
    fn (data: bytes, level: i32) -> bytes
    + compresses data using LZ77 + Huffman at the given level (0-9)
    + level 0 emits stored blocks only
    # compression
  zlib.inflate
    fn (data: bytes) -> result[bytes, string]
    + decompresses a raw deflate stream
    - returns error on truncated input or invalid block type
    # decompression
  zlib.adler32
    fn (data: bytes) -> u32
    + returns the Adler-32 checksum
    # checksum
  zlib.crc32
    fn (data: bytes) -> u32
    + returns the CRC-32 checksum using the IEEE polynomial
    # checksum
  zlib.zlib_wrap
    fn (data: bytes, level: i32) -> bytes
    + wraps deflate output with a 2-byte zlib header and trailing Adler-32
    # framing
    -> zlib.deflate
    -> zlib.adler32
  zlib.zlib_unwrap
    fn (data: bytes) -> result[bytes, string]
    + validates the zlib header and Adler-32 trailer, returning the inflated payload
    - returns error on header or checksum mismatch
    # framing
    -> zlib.inflate
    -> zlib.adler32
  zlib.gzip_wrap
    fn (data: bytes, level: i32, filename: optional[string]) -> bytes
    + wraps deflate output with a gzip header, optional filename field, and trailing CRC-32 and original size
    # framing
    -> zlib.deflate
    -> zlib.crc32
  zlib.gzip_unwrap
    fn (data: bytes) -> result[bytes, string]
    + parses the gzip header, inflates the payload, and validates CRC-32 and size
    - returns error on magic number or checksum mismatch
    # framing
    -> zlib.inflate
    -> zlib.crc32
