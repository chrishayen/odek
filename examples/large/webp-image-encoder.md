# Requirement: "a WebP image encoder"

A lossless WebP encoder. Real work lives in std (bitstream writing, predictive transforms, entropy coding); the project layer orchestrates header, image data, and container assembly.

std
  std.bitstream
    std.bitstream.new_writer
      fn () -> bit_writer
      + creates an empty bit writer
      # bitstream
    std.bitstream.write_bits
      fn (w: bit_writer, value: u32, n_bits: i32) -> bit_writer
      + appends the low n_bits of value, LSB first
      # bitstream
    std.bitstream.finish
      fn (w: bit_writer) -> bytes
      + returns the byte-aligned buffer, padding the final byte with zeros
      # bitstream
  std.compression
    std.compression.huffman_build
      fn (frequencies: list[i32]) -> huffman_table
      + builds a canonical Huffman table for the given symbol frequencies
      + produces equal-length codes when all frequencies are equal
      # compression
    std.compression.huffman_encode_symbol
      fn (table: huffman_table, symbol: i32) -> tuple[u32, i32]
      + returns the code and its bit length for a symbol
      # compression
    std.compression.lz77_match
      fn (data: bytes, at: i32, window: i32) -> optional[lz_match]
      + returns the longest back-reference within window bytes
      # compression
  std.image
    std.image.pixels_from_rgba
      fn (data: bytes, width: i32, height: i32) -> result[pixel_buffer, string]
      + validates that data length equals width*height*4
      - returns error on mismatched dimensions
      # image
  std.io
    std.io.write_u32_le
      fn (out: bytes, value: u32) -> bytes
      + appends a little-endian u32
      # io

webp
  webp.encode_lossless
    fn (pixels: bytes, width: i32, height: i32) -> result[bytes, string]
    + returns a complete RIFF/WebP file with a lossless VP8L payload
    - returns error when width or height is zero or exceeds 16383
    # encoding
    -> std.image.pixels_from_rgba
    -> std.bitstream.new_writer
    -> std.bitstream.write_bits
    -> std.bitstream.finish
  webp.apply_color_transform
    fn (pixels: pixel_buffer) -> pixel_buffer
    + subtracts a per-block green channel predictor from red and blue channels
    # transform
  webp.apply_predictor_transform
    fn (pixels: pixel_buffer) -> pixel_buffer
    + replaces each pixel with its residual against a chosen predictor mode
    # transform
  webp.build_huffman_codes
    fn (pixels: pixel_buffer) -> huffman_table
    + builds Huffman codes over the transformed pixel stream
    # entropy
    -> std.compression.huffman_build
  webp.write_vp8l_payload
    fn (pixels: pixel_buffer, table: huffman_table) -> bytes
    + emits the VP8L signature, transform flags, and entropy-coded pixel stream
    # payload
    -> std.compression.lz77_match
    -> std.compression.huffman_encode_symbol
    -> std.bitstream.write_bits
  webp.wrap_riff
    fn (payload: bytes) -> bytes
    + wraps the payload in a RIFF "WEBPVP8L" container with the correct chunk sizes
    # container
    -> std.io.write_u32_le
