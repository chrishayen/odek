# Requirement: "a JSON object signing and encryption library"

Implements JWS (signing) and JWE (encryption) envelopes over JSON payloads. The project surface is small; std carries the cryptographic primitives.

std
  std.encoding
    std.encoding.base64url_encode
      @ (data: bytes) -> string
      + encodes bytes to base64url without padding
      # encoding
    std.encoding.base64url_decode
      @ (encoded: string) -> result[bytes, string]
      + decodes base64url with or without padding
      - returns error on characters outside the base64url alphabet
      # encoding
  std.crypto
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256 of data under key
      + returns 32 bytes
      # cryptography
    std.crypto.aes_gcm_encrypt
      @ (key: bytes, iv: bytes, plaintext: bytes, aad: bytes) -> tuple[bytes, bytes]
      + returns (ciphertext, auth_tag) for AES-GCM
      # cryptography
    std.crypto.aes_gcm_decrypt
      @ (key: bytes, iv: bytes, ciphertext: bytes, tag: bytes, aad: bytes) -> result[bytes, string]
      + returns plaintext on success
      - returns error when the tag does not verify
      # cryptography
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically strong random bytes
      # cryptography
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a flat string-to-string map as JSON
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a flat JSON object
      - returns error on invalid JSON or non-object root
      # serialization

jose
  jose.sign_hs256
    @ (payload: map[string, string], secret: bytes) -> result[string, string]
    + returns a compact JWS string "header.payload.signature"
    - returns error when secret is empty
    # signing
    -> std.json.encode_object
    -> std.encoding.base64url_encode
    -> std.crypto.hmac_sha256
  jose.verify_hs256
    @ (token: string, secret: bytes) -> result[map[string, string], string]
    + returns the payload on valid signature
    - returns error on malformed segments or bad signature
    # verification
    -> std.encoding.base64url_decode
    -> std.crypto.hmac_sha256
    -> std.json.parse_object
  jose.encrypt_a256gcm
    @ (payload: map[string, string], key: bytes) -> result[string, string]
    + returns a compact JWE "header.empty_key.iv.ciphertext.tag"
    - returns error when key is not 32 bytes
    # encryption
    -> std.crypto.random_bytes
    -> std.crypto.aes_gcm_encrypt
    -> std.encoding.base64url_encode
    -> std.json.encode_object
  jose.decrypt_a256gcm
    @ (token: string, key: bytes) -> result[map[string, string], string]
    + returns the payload on successful authenticated decryption
    - returns error when tag fails to verify
    - returns error when the token does not have exactly five segments
    # decryption
    -> std.encoding.base64url_decode
    -> std.crypto.aes_gcm_decrypt
    -> std.json.parse_object
