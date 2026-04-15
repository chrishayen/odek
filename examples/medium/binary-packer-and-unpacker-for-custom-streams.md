# Requirement: "a binary packer and unpacker for custom binary streams"

A fluent-ish builder that appends fixed-width integers, bytes, and length-prefixed strings, and a matching reader.

std: (all units exist)

binpacker
  binpacker.new_packer
    fn () -> packer_state
    + creates an empty packer
    # construction
  binpacker.push_u8
    fn (state: packer_state, value: u8) -> packer_state
    + appends one byte
    # write
  binpacker.push_u16_be
    fn (state: packer_state, value: u16) -> packer_state
    + appends a big-endian u16
    # write
  binpacker.push_u32_be
    fn (state: packer_state, value: u32) -> packer_state
    + appends a big-endian u32
    # write
  binpacker.push_u64_be
    fn (state: packer_state, value: u64) -> packer_state
    + appends a big-endian u64
    # write
  binpacker.push_bytes
    fn (state: packer_state, data: bytes) -> packer_state
    + appends the raw bytes with no length prefix
    # write
  binpacker.push_length_string
    fn (state: packer_state, s: string) -> packer_state
    + appends a u32 big-endian length followed by the utf-8 bytes
    # write
  binpacker.to_bytes
    fn (state: packer_state) -> bytes
    + returns the accumulated bytes
    # finalization
  binpacker.new_unpacker
    fn (data: bytes) -> unpacker_state
    + creates an unpacker positioned at offset 0
    # construction
  binpacker.read_u16_be
    fn (state: unpacker_state) -> result[tuple[u16, unpacker_state], string]
    + reads a big-endian u16 and advances the cursor
    - returns error on end of buffer
    # read
  binpacker.read_u32_be
    fn (state: unpacker_state) -> result[tuple[u32, unpacker_state], string]
    + reads a big-endian u32 and advances the cursor
    - returns error on end of buffer
    # read
  binpacker.read_length_string
    fn (state: unpacker_state) -> result[tuple[string, unpacker_state], string]
    + reads a u32 length then that many bytes as a utf-8 string
    - returns error when the declared length exceeds the remaining buffer
    # read
