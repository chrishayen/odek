# Requirement: "a brotli decompressor"

Decompresses brotli-compressed byte streams. The state machine over the bit stream lives in project code; bit-level reading lives in std.

std
  std.bits
    std.bits.new_reader
      @ (data: bytes) -> bit_reader
      + creates a little-endian bit reader positioned at bit 0
      # bit_reading
    std.bits.read_bits
      @ (r: bit_reader, n: i32) -> result[tuple[u32, bit_reader], string]
      + reads n bits (1..24) and advances the reader
      - returns error when reader is exhausted
      # bit_reading

brotli
  brotli.decompress
    @ (compressed: bytes) -> result[bytes, string]
    + decompresses a complete brotli stream to its original bytes
    - returns error on truncated input
    - returns error on invalid header
    # decompression
    -> std.bits.new_reader
  brotli.read_header
    @ (r: bit_reader) -> result[tuple[brotli_header, bit_reader], string]
    + parses the window size and metadata flags from the stream header
    - returns error on reserved bit set
    # header_parsing
    -> std.bits.read_bits
  brotli.decode_meta_block
    @ (r: bit_reader, out: bytes) -> result[tuple[bytes, bit_reader, bool], string]
    + decodes one meta-block and returns updated output and whether this is the final block
    - returns error on invalid huffman tree
    # block_decoding
    -> std.bits.read_bits
  brotli.build_huffman_tree
    @ (code_lengths: list[i32]) -> result[huffman_tree, string]
    + constructs a canonical huffman tree from code lengths
    - returns error when lengths are not a valid prefix code
    # huffman
