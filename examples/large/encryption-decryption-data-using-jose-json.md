# Requirement: "JOSE encryption and decryption of data"

Implements the JWE compact serialization with symmetric and asymmetric key-wrap modes.

std
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
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.crypto
    std.crypto.aes_gcm_seal
      @ (key: bytes, nonce: bytes, aad: bytes, plaintext: bytes) -> tuple[bytes, bytes]
      + returns (ciphertext, 16-byte tag)
      # cryptography
    std.crypto.aes_gcm_open
      @ (key: bytes, nonce: bytes, aad: bytes, ciphertext: bytes, tag: bytes) -> result[bytes, string]
      + returns the plaintext after verifying the tag
      - returns error when the tag does not match
      # cryptography
    std.crypto.aes_key_wrap
      @ (kek: bytes, plaintext_key: bytes) -> result[bytes, string]
      + returns RFC 3394 key-wrapped output
      - returns error when kek length is not 16, 24, or 32
      # cryptography
    std.crypto.aes_key_unwrap
      @ (kek: bytes, wrapped: bytes) -> result[bytes, string]
      + returns the unwrapped key
      - returns error when the integrity check fails
      # cryptography
    std.crypto.rsa_encrypt_oaep
      @ (public_key: bytes, plaintext: bytes) -> result[bytes, string]
      + returns RSA-OAEP ciphertext
      # cryptography
    std.crypto.rsa_decrypt_oaep
      @ (private_key: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns the plaintext
      - returns error when ciphertext is malformed
      # cryptography
    std.crypto.random_bytes
      @ (length: i32) -> bytes
      + returns cryptographically random bytes
      # randomness

jose
  jose.encrypt_compact
    @ (header: map[string, string], recipient_key: bytes, plaintext: bytes) -> result[string, string]
    + returns a compact JWE "h.ek.iv.ct.tag" string
    - returns error when header is missing alg or enc
    - returns error when alg is not recognized
    # encryption
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.random_bytes
    -> std.crypto.aes_gcm_seal
  jose.decrypt_compact
    @ (compact: string, recipient_key: bytes) -> result[bytes, string]
    + returns the plaintext when all authentication checks pass
    - returns error when the JWE does not have exactly five segments
    - returns error when tag verification fails
    # decryption
    -> std.encoding.base64url_decode
    -> std.json.parse_object
    -> std.crypto.aes_gcm_open
  jose.wrap_content_key
    @ (alg: string, recipient_key: bytes, cek: bytes) -> result[bytes, string]
    + returns the encrypted content encryption key for the chosen alg
    - returns error when alg is unknown
    # key_management
    -> std.crypto.aes_key_wrap
    -> std.crypto.rsa_encrypt_oaep
  jose.unwrap_content_key
    @ (alg: string, recipient_key: bytes, wrapped: bytes) -> result[bytes, string]
    + returns the content encryption key
    - returns error when alg is unknown
    - returns error when unwrap integrity fails
    # key_management
    -> std.crypto.aes_key_unwrap
    -> std.crypto.rsa_decrypt_oaep
  jose.build_protected_header
    @ (alg: string, enc: string, extras: map[string, string]) -> string
    + returns the base64url-encoded protected header JSON
    # encryption
    -> std.json.encode_object
    -> std.encoding.base64url_encode
