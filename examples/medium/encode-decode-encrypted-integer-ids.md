# Requirement: "encode and decode encrypted integer IDs"

Integers are encrypted with a symmetric key and serialized to a short URL-safe string, so database ids can be exposed to clients without leaking sequence order.

std
  std.crypto
    std.crypto.aes_ctr_encrypt
      @ (key: bytes, nonce: bytes, plaintext: bytes) -> result[bytes, string]
      + returns the AES-CTR ciphertext
      - returns error when key length is not 16, 24, or 32
      # cryptography
    std.crypto.aes_ctr_decrypt
      @ (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns the AES-CTR plaintext
      - returns error when key length is not 16, 24, or 32
      # cryptography
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256 of data under key
      + returns 32 bytes
      # cryptography
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
    std.encoding.base64url_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes base64url with or without padding
      - returns error on characters outside the alphabet
      # encoding
    std.encoding.encode_u64_be
      @ (value: u64) -> bytes
      + encodes a u64 as 8 big-endian bytes
      # encoding
    std.encoding.decode_u64_be
      @ (data: bytes) -> result[u64, string]
      + decodes 8 big-endian bytes into a u64
      - returns error when input is shorter than 8 bytes
      # encoding

encid
  encid.encode
    @ (id: u64, key: bytes) -> result[string, string]
    + returns an opaque base64url string representing the encrypted id
    - returns error when key is not 16, 24, or 32 bytes
    ? nonce is derived from an HMAC of the id so the same id always produces the same ciphertext
    # encoding
    -> std.encoding.encode_u64_be
    -> std.crypto.hmac_sha256
    -> std.crypto.aes_ctr_encrypt
    -> std.encoding.base64url_encode
  encid.decode
    @ (encoded: string, key: bytes) -> result[u64, string]
    + returns the original integer id
    - returns error when the input is not valid base64url
    - returns error when the ciphertext length is unexpected
    - returns error when the nonce does not match the recomputed HMAC
    # decoding
    -> std.encoding.base64url_decode
    -> std.crypto.aes_ctr_decrypt
    -> std.crypto.hmac_sha256
    -> std.encoding.decode_u64_be
