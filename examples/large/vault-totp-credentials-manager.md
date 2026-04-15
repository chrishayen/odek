# Requirement: "a credentials manager with TOTP 2FA support"

Stores credentials encrypted at rest and generates TOTP codes for 2FA-enabled entries.

std
  std.crypto
    std.crypto.hmac_sha1
      fn (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA1 of data under key
      + returns 20 bytes
      # cryptography
    std.crypto.aes256_gcm_encrypt
      fn (key: bytes, nonce: bytes, plaintext: bytes) -> result[bytes, string]
      + encrypts plaintext with AES-256-GCM
      - returns error when key is not 32 bytes
      # cryptography
    std.crypto.aes256_gcm_decrypt
      fn (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      + decrypts AES-256-GCM ciphertext
      - returns error on authentication failure
      # cryptography
    std.crypto.pbkdf2_sha256
      fn (password: string, salt: bytes, iterations: i32, length: i32) -> bytes
      + derives a key of the requested length from the password
      # key_derivation
    std.crypto.random_bytes
      fn (length: i32) -> bytes
      + returns cryptographically random bytes
      # randomness
  std.encoding
    std.encoding.base32_decode
      fn (encoded: string) -> result[bytes, string]
      + decodes base32 input
      - returns error on invalid characters
      # encoding
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads file bytes
      - returns error when missing
      # fs
    std.fs.write_all
      fn (path: string, contents: bytes) -> result[void, string]
      + writes bytes to path
      - returns error on I/O failure
      # fs

vault
  vault.create
    fn (master_password: string, path: string) -> result[void, string]
    + creates an empty encrypted vault file secured by the master password
    - returns error when path already exists
    # vault_init
    -> std.crypto.random_bytes
    -> std.crypto.pbkdf2_sha256
    -> std.crypto.aes256_gcm_encrypt
    -> std.fs.write_all
  vault.open
    fn (master_password: string, path: string) -> result[vault_state, string]
    + unlocks an existing vault and returns its in-memory state
    - returns error when the password is incorrect
    - returns error when the vault file is corrupt
    # vault_init
    -> std.fs.read_all
    -> std.crypto.pbkdf2_sha256
    -> std.crypto.aes256_gcm_decrypt
  vault.save
    fn (state: vault_state, path: string) -> result[void, string]
    + serializes the vault and writes it back encrypted
    - returns error on I/O failure
    # persistence
    -> std.crypto.random_bytes
    -> std.crypto.aes256_gcm_encrypt
    -> std.fs.write_all
  vault.add_entry
    fn (state: vault_state, name: string, username: string, password: string) -> result[vault_state, string]
    + stores a new credential under the given name
    - returns error when an entry with that name already exists
    # credentials
  vault.get_entry
    fn (state: vault_state, name: string) -> result[credential, string]
    + returns the credential for the given name
    - returns error when name is not found
    # credentials
  vault.remove_entry
    fn (state: vault_state, name: string) -> result[vault_state, string]
    + deletes the named entry
    - returns error when name is not found
    # credentials
  vault.set_totp_secret
    fn (state: vault_state, name: string, secret_base32: string) -> result[vault_state, string]
    + attaches a TOTP shared secret to an existing entry
    - returns error when secret is not valid base32
    # two_factor
    -> std.encoding.base32_decode
  vault.generate_totp
    fn (state: vault_state, name: string) -> result[string, string]
    + returns the current 6-digit TOTP code for the entry
    - returns error when the entry has no TOTP secret
    # two_factor
    -> std.time.now_seconds
    -> std.crypto.hmac_sha1
  vault.change_master_password
    fn (state: vault_state, old: string, new: string) -> result[vault_state, string]
    + rotates the master password and re-derives the data key
    - returns error when old password does not match
    # vault_admin
    -> std.crypto.pbkdf2_sha256
