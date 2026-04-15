# Requirement: "an LSB steganography library that hides a file inside an image"

Embeds arbitrary payload bytes into the least significant bits of a raw RGB pixel buffer, and extracts them back.

std
  std.image
    std.image.decode_rgb
      fn (raw: bytes) -> result[list[u8], string]
      + decodes an image into a flat RGB pixel buffer (one byte per channel)
      - returns error on unsupported formats
      # image_codec
    std.image.encode_rgb
      fn (pixels: list[u8], width: i32, height: i32) -> bytes
      + encodes a flat RGB buffer back into an image container
      # image_codec

stego
  stego.capacity_bytes
    fn (pixel_count: i32) -> i32
    + returns the maximum payload size (in bytes) a pixel buffer can hold
    ? each pixel channel stores one bit of payload
    # capacity
  stego.embed
    fn (pixels: list[u8], payload: bytes) -> result[list[u8], string]
    + returns a new pixel buffer with the payload hidden in the low bits
    + prefixes the embedded data with a 32-bit length header
    - returns error when payload exceeds capacity
    # embedding
  stego.extract
    fn (pixels: list[u8]) -> result[bytes, string]
    + reads the length header and returns the original payload
    - returns error when the header claims a length larger than capacity
    # extraction
