# Requirement: "a zero-copy memory manipulation library"

Reinterpret byte buffers as typed values without allocation or copying. Project surface stays small; alignment and bounds checks are the real work.

std
  std.mem
    std.mem.align_of
      fn (type_name: string) -> u32
      + returns the required alignment in bytes for a known primitive type name
      + returns 1 for "u8" and "i8"
      - returns 0 for unknown type names
      # memory
    std.mem.size_of
      fn (type_name: string) -> u32
      + returns the byte width of a known primitive type name
      - returns 0 for unknown type names
      # memory

zerocopy
  zerocopy.view_u32_le
    fn (buf: bytes, offset: u32) -> result[u32, string]
    + returns the little-endian u32 at the given offset without copying
    - returns error when offset + 4 exceeds buffer length
    - returns error when offset is not 4-byte aligned
    # reinterpret
    -> std.mem.align_of
    -> std.mem.size_of
  zerocopy.view_i64_le
    fn (buf: bytes, offset: u32) -> result[i64, string]
    + returns the little-endian i64 at the given offset without copying
    - returns error when offset + 8 exceeds buffer length
    - returns error when offset is not 8-byte aligned
    # reinterpret
    -> std.mem.align_of
    -> std.mem.size_of
  zerocopy.slice
    fn (buf: bytes, offset: u32, length: u32) -> result[bytes, string]
    + returns a subrange of the buffer sharing the same backing storage
    - returns error when offset + length exceeds buffer length
    ? the returned bytes alias the input; callers must not mutate either while the view is live
    # slicing
  zerocopy.write_u32_le
    fn (buf: bytes, offset: u32, value: u32) -> result[void, string]
    + writes value in little-endian form at offset, mutating the buffer in place
    - returns error when offset + 4 exceeds buffer length
    - returns error when offset is not 4-byte aligned
    # in_place_write
    -> std.mem.align_of
