# Requirement: "a secrets manager that encrypts and decrypts using ssh keys"

Encrypts a secret payload to the holder of an ssh public key and decrypts it with the matching private key.

std
  std.fs
    std.fs.read_all_bytes
      fn (path: string) -> result[bytes, string]
      + reads the entire file as bytes
      # filesystem
  std.ssh_key
    std.ssh_key.parse_public
      fn (raw: bytes) -> result[ssh_public_key, string]
      + parses an openssh-format public key
      - returns error on malformed input
      # keys
    std.ssh_key.parse_private
      fn (raw: bytes, passphrase: optional[string]) -> result[ssh_private_key, string]
      + parses a private key, optionally decrypting with the passphrase
      - returns error on wrong passphrase or unsupported key type
      # keys
  std.crypto
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
    std.crypto.rsa_encrypt
      fn (key: ssh_public_key, plaintext: bytes) -> result[bytes, string]
      + encrypts plaintext with an rsa public key using oaep padding
      - returns error when the key is not rsa
      # cryptography
    std.crypto.rsa_decrypt
      fn (key: ssh_private_key, ciphertext: bytes) -> result[bytes, string]
      + decrypts ciphertext with the matching rsa private key
      - returns error on wrong key or corrupt ciphertext
      # cryptography
    std.crypto.aes_gcm_seal
      fn (key: bytes, nonce: bytes, plaintext: bytes) -> result[bytes, string]
      + seals plaintext under a 32-byte key and 12-byte nonce
      # cryptography
    std.crypto.aes_gcm_open
      fn (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      + opens sealed ciphertext
      - returns error on authentication failure
      # cryptography

secrets
  secrets.encrypt
    fn (public_key_path: string, plaintext: bytes) -> result[bytes, string]
    + generates a random symmetric key, seals the payload with aes-gcm, wraps the key with the ssh public key, and packs the result
    - returns error when the key cannot be read or is not an rsa key
    # encryption
    -> std.fs.read_all_bytes
    -> std.ssh_key.parse_public
    -> std.crypto.random_bytes
    -> std.crypto.aes_gcm_seal
    -> std.crypto.rsa_encrypt
  secrets.decrypt
    fn (private_key_path: string, passphrase: optional[string], vault_bytes: bytes) -> result[bytes, string]
    + unwraps the symmetric key with the private key and opens the payload
    - returns error when the passphrase is wrong or the vault is corrupt
    # decryption
    -> std.fs.read_all_bytes
    -> std.ssh_key.parse_private
    -> std.crypto.rsa_decrypt
    -> std.crypto.aes_gcm_open
