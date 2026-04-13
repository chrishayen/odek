# Requirement: "a byte order encoder and decoder for big-endian, little-endian, and native orders"

Pure functions over byte buffers. No std primitives needed beyond what the host provides.

std: (all units exist)

byteorder
  byteorder.write_u32_be
    @ (value: u32) -> bytes
    + returns 4 bytes in big-endian order
    # encoding
  byteorder.write_u32_le
    @ (value: u32) -> bytes
    + returns 4 bytes in little-endian order
    # encoding
  byteorder.read_u32_be
    @ (buf: bytes) -> result[u32, string]
    + decodes the first 4 bytes as big-endian
    - returns error when the buffer is shorter than 4 bytes
    # decoding
  byteorder.read_u32_le
    @ (buf: bytes) -> result[u32, string]
    + decodes the first 4 bytes as little-endian
    - returns error when the buffer is shorter than 4 bytes
    # decoding
  byteorder.write_u64_be
    @ (value: u64) -> bytes
    + returns 8 bytes in big-endian order
    # encoding
  byteorder.read_u64_be
    @ (buf: bytes) -> result[u64, string]
    + decodes the first 8 bytes as big-endian
    - returns error when the buffer is shorter than 8 bytes
    # decoding
