# Requirement: "a perceptual image hashing library"

Computes average, difference, and perceptual hashes so near-duplicate images compare as short fingerprints.

std
  std.image
    std.image.decode
      @ (raw: bytes) -> result[image, string]
      + decodes a bitmap image into width, height, and pixel data
      - returns error on unsupported or malformed input
      # image
    std.image.to_grayscale
      @ (img: image) -> image
      + returns a single-channel grayscale copy
      # image
    std.image.resize
      @ (img: image, width: i32, height: i32) -> image
      + returns a resized copy using bilinear sampling
      # image
  std.math
    std.math.dct2
      @ (data: list[f64], n: i32) -> list[f64]
      + returns the two-dimensional discrete cosine transform of an n-by-n block
      # math

imghash
  imghash.average_hash
    @ (raw: bytes) -> result[u64, string]
    + returns a 64-bit hash where each bit is 1 when the pixel exceeds the 8x8 mean
    - returns error on decode failure
    # hashing
    -> std.image.decode
    -> std.image.to_grayscale
    -> std.image.resize
  imghash.difference_hash
    @ (raw: bytes) -> result[u64, string]
    + returns a 64-bit hash by comparing adjacent pixels in a 9x8 resized image
    - returns error on decode failure
    # hashing
    -> std.image.decode
    -> std.image.to_grayscale
    -> std.image.resize
  imghash.perceptual_hash
    @ (raw: bytes) -> result[u64, string]
    + resizes to 32x32, applies a 2D DCT, and hashes the top-left 8x8 coefficients against their median
    - returns error on decode failure
    # hashing
    -> std.image.decode
    -> std.image.to_grayscale
    -> std.image.resize
    -> std.math.dct2
  imghash.hamming_distance
    @ (a: u64, b: u64) -> i32
    + returns the number of differing bits between two hashes
    # comparison
