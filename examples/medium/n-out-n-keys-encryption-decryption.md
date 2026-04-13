# Requirement: "an n-out-of-N keys encryption and decryption framework based on Shamir's Secret Sharing"

A data encryption key is split into N shares via Shamir's scheme; any threshold t of them can reconstruct the key to decrypt the payload.

std
  std.crypto
    std.crypto.random_bytes
      @ (count: i32) -> bytes
      + returns count cryptographically random bytes
      # cryptography
    std.crypto.aes_gcm_seal
      @ (key: bytes, nonce: bytes, plaintext: bytes) -> bytes
      + returns ciphertext with appended authentication tag
      # cryptography
    std.crypto.aes_gcm_open
      @ (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns plaintext when the tag verifies
      - returns error when the tag does not verify
      # cryptography
  std.math
    std.math.gf256_eval_poly
      @ (coeffs: list[u8], x: u8) -> u8
      + evaluates a polynomial over GF(2^8) at x using Horner's rule
      # finite_field
    std.math.gf256_lagrange_interpolate_at_zero
      @ (xs: list[u8], ys: list[u8]) -> result[u8, string]
      + returns the constant term via Lagrange interpolation at x=0 over GF(2^8)
      - returns error when xs contains duplicates or lengths differ
      # finite_field

secret_split
  secret_split.split_key
    @ (key: bytes, threshold: i32, shares: i32) -> result[list[tuple[u8, bytes]], string]
    + returns shares pairs of (x, share_bytes) where any threshold of them can reconstruct key
    - returns error when threshold < 2, shares < threshold, or shares > 255
    # sharing
    -> std.crypto.random_bytes
    -> std.math.gf256_eval_poly
  secret_split.combine_key
    @ (parts: list[tuple[u8, bytes]]) -> result[bytes, string]
    + reconstructs the original key byte-by-byte from threshold shares
    - returns error when fewer than 2 parts are given or share lengths disagree
    # reconstruction
    -> std.math.gf256_lagrange_interpolate_at_zero
  secret_split.encrypt
    @ (plaintext: bytes, threshold: i32, shares: i32) -> result[tuple[bytes, list[tuple[u8, bytes]]], string]
    + generates a fresh data key, seals the plaintext, and returns the ciphertext and shares
    - returns error when threshold/shares are invalid
    # encryption
    -> std.crypto.random_bytes
    -> std.crypto.aes_gcm_seal
  secret_split.decrypt
    @ (ciphertext: bytes, parts: list[tuple[u8, bytes]]) -> result[bytes, string]
    + combines the shares and opens the ciphertext with the reconstructed key
    - returns error when shares are insufficient or authentication fails
    # decryption
    -> std.crypto.aes_gcm_open
