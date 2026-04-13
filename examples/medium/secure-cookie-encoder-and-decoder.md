# Requirement: "a secure cookie encoder and decoder"

Serializes cookie values, encrypts them, and verifies integrity on decode.

std
  std.crypto
    std.crypto.aes256_gcm_encrypt
      @ (key: bytes, nonce: bytes, plaintext: bytes) -> result[bytes, string]
      + encrypts plaintext with AES-256-GCM
      - returns error when key is not 32 bytes
      # cryptography
    std.crypto.aes256_gcm_decrypt
      @ (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      + decrypts AES-256-GCM ciphertext
      - returns error on authentication failure
      # cryptography
    std.crypto.random_bytes
      @ (length: i32) -> bytes
      + returns cryptographically random bytes
      # randomness
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
    std.encoding.base64url_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes base64url with or without padding
      - returns error on invalid characters
      # encoding
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

securecookie
  securecookie.new
    @ (key: bytes, max_age_seconds: i64) -> cookie_codec
    + creates a codec bound to the encryption key and max age
    - returns error when key is not 32 bytes
    # construction
  securecookie.encode
    @ (codec: cookie_codec, name: string, value: string) -> result[string, string]
    + returns an encrypted, timestamped, base64url-encoded cookie value
    - returns error when encryption fails
    # encoding
    -> std.crypto.random_bytes
    -> std.time.now_seconds
    -> std.crypto.aes256_gcm_encrypt
    -> std.encoding.base64url_encode
  securecookie.decode
    @ (codec: cookie_codec, name: string, encoded: string) -> result[string, string]
    + returns the original value when the cookie is valid and fresh
    - returns error when the MAC fails
    - returns error when the cookie exceeds max age
    # decoding
    -> std.encoding.base64url_decode
    -> std.crypto.aes256_gcm_decrypt
    -> std.time.now_seconds
