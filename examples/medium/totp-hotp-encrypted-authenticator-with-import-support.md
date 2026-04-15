# Requirement: "an encrypted TOTP/HOTP authenticator library with import support"

Generates one-time codes, stores accounts under a passphrase-derived key, and accepts imports from an otpauth-URI list.

std
  std.crypto
    std.crypto.hmac_sha1
      fn (key: bytes, data: bytes) -> bytes
      + returns HMAC-SHA1 of data under key
      + result is 20 bytes
      # cryptography
    std.crypto.pbkdf2_sha256
      fn (password: bytes, salt: bytes, iterations: i32, length: i32) -> bytes
      + derives a key of the requested length
      + returns identical output for identical inputs
      # key_derivation
    std.crypto.aead_encrypt
      fn (key: bytes, nonce: bytes, plaintext: bytes) -> bytes
      + returns ciphertext with authentication tag appended
      # authenticated_encryption
    std.crypto.aead_decrypt
      fn (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns plaintext when tag verifies
      - returns error on tag mismatch
      # authenticated_encryption
  std.encoding
    std.encoding.base32_decode
      fn (s: string) -> result[bytes, string]
      + decodes RFC 4648 base32, ignoring padding and case
      - returns error on invalid characters
      # encoding
  std.url
    std.url.parse_otpauth
      fn (uri: string) -> result[map[string, string], string]
      + extracts type, issuer, account, secret, algorithm, digits, period
      - returns error when scheme is not otpauth
      # uri_parsing

authenticator
  authenticator.hotp
    fn (secret: bytes, counter: u64, digits: i32) -> string
    + returns a zero-padded decimal code of the requested length
    ? uses the dynamic-truncation algorithm
    # hotp
    -> std.crypto.hmac_sha1
  authenticator.totp
    fn (secret: bytes, unix_time: i64, period: i32, digits: i32) -> string
    + returns the code for the time window containing unix_time
    # totp
  authenticator.import_uri_list
    fn (text: string) -> result[list[account], string]
    + parses one otpauth URI per line into account records
    - returns error on any malformed line
    # import
    -> std.url.parse_otpauth
    -> std.encoding.base32_decode
  authenticator.seal_vault
    fn (accounts: list[account], passphrase: string, salt: bytes, nonce: bytes) -> bytes
    + returns a passphrase-encrypted blob of the account list
    # vault_seal
    -> std.crypto.pbkdf2_sha256
    -> std.crypto.aead_encrypt
  authenticator.open_vault
    fn (blob: bytes, passphrase: string, salt: bytes, nonce: bytes) -> result[list[account], string]
    + returns the account list when passphrase is correct
    - returns error on wrong passphrase or tampered blob
    # vault_open
    -> std.crypto.pbkdf2_sha256
    -> std.crypto.aead_decrypt
