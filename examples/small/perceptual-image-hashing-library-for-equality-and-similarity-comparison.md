# Requirement: "a perceptual image hashing library for equality and similarity comparison"

Produces compact fingerprints of images so visually similar images hash to similar values; similarity is measured with Hamming distance.

std: (all units exist)

phash
  phash.compute
    fn (pixels: list[u8], width: i32, height: i32) -> bytes
    + returns a 64-bit perceptual hash computed from a DCT of a downscaled grayscale image
    ? accepts interleaved RGB or RGBA; channel count inferred from pixels.len / (width * height)
    # hashing
  phash.distance
    fn (a: bytes, b: bytes) -> i32
    + returns the Hamming distance (number of differing bits) between two hashes
    - returns -1 when the hashes differ in length
    # comparison
  phash.equal
    fn (a: bytes, b: bytes) -> bool
    + returns true when the hashes are byte-for-byte identical
    # comparison
  phash.similar
    fn (a: bytes, b: bytes, threshold: i32) -> bool
    + returns true when the Hamming distance is at most threshold
    # comparison
