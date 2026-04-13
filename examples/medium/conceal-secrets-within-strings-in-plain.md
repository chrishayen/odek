# Requirement: "a library that conceals secrets within innocuous strings using invisible characters"

Encrypts a short secret under a passphrase, encodes the ciphertext as zero-width unicode characters, and hides the result inside a cover string.

std
  std.crypto
    std.crypto.derive_key
      @ (passphrase: string, salt: bytes) -> bytes
      + returns a 32-byte symmetric key
      # cryptography
    std.crypto.encrypt_aead
      @ (key: bytes, plaintext: bytes, nonce: bytes) -> bytes
      + returns ciphertext with authentication tag appended
      # cryptography
    std.crypto.decrypt_aead
      @ (key: bytes, ciphertext: bytes, nonce: bytes) -> result[bytes, string]
      + returns the plaintext when the tag is valid
      - returns error on tampering
      # cryptography
  std.random
    std.random.bytes
      @ (n: u32) -> bytes
      + returns n cryptographically random bytes
      # randomness

steganographic_text
  steganographic_text.encode_invisible
    @ (data: bytes) -> string
    + maps bytes onto a fixed alphabet of zero-width unicode characters
    # encoding
  steganographic_text.decode_invisible
    @ (encoded: string) -> result[bytes, string]
    + recovers the original bytes from the zero-width sequence
    - returns error when a character is outside the alphabet
    # decoding
  steganographic_text.hide
    @ (cover: string, secret: string, passphrase: string) -> string
    + returns the cover text with the encrypted secret woven between its characters
    # concealment
    -> std.random.bytes
    -> std.crypto.derive_key
    -> std.crypto.encrypt_aead
  steganographic_text.reveal
    @ (carrier: string, passphrase: string) -> result[string, string]
    + extracts and decrypts the hidden payload
    - returns error when the passphrase does not match
    - returns error when no hidden payload is present
    # extraction
    -> std.crypto.derive_key
    -> std.crypto.decrypt_aead
