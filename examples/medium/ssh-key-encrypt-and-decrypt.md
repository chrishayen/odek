# Requirement: "encrypt and decrypt data using SSH keys"

Uses the public half of an SSH key pair as a recipient and the private half to unseal. std exposes the needed asymmetric and symmetric primitives.

std
  std.crypto
    std.crypto.rsa_encrypt_oaep
      @ (public_key: bytes, plaintext: bytes) -> result[bytes, string]
      + returns RSA-OAEP ciphertext
      - returns error when plaintext is longer than the key allows
      # cryptography
    std.crypto.rsa_decrypt_oaep
      @ (private_key: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns the decrypted plaintext
      - returns error when the ciphertext is malformed
      # cryptography
    std.crypto.aes_gcm_seal
      @ (key: bytes, nonce: bytes, plaintext: bytes) -> bytes
      + returns ciphertext with appended authentication tag
      # cryptography
    std.crypto.aes_gcm_open
      @ (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns plaintext after verifying the tag
      - returns error when the authentication tag does not match
      # cryptography
    std.crypto.random_bytes
      @ (length: i32) -> bytes
      + returns cryptographically random bytes
      # randomness
  std.ssh
    std.ssh.parse_public_key
      @ (authorized_line: string) -> result[ssh_public_key, string]
      + parses an "ssh-rsa AAAA..." line into a typed key
      - returns error on unsupported key types
      - returns error on malformed base64 body
      # ssh
    std.ssh.parse_private_key
      @ (pem: string, passphrase: optional[string]) -> result[ssh_private_key, string]
      + parses an OpenSSH PEM private key
      - returns error on wrong passphrase
      # ssh

ssh_crypt
  ssh_crypt.seal
    @ (public_key_line: string, plaintext: bytes) -> result[bytes, string]
    + returns a self-contained sealed blob: RSA-OAEP-wrapped content key plus AES-GCM body
    - returns error when the key line does not parse
    # sealing
    -> std.ssh.parse_public_key
    -> std.crypto.random_bytes
    -> std.crypto.rsa_encrypt_oaep
    -> std.crypto.aes_gcm_seal
  ssh_crypt.open
    @ (private_key_pem: string, passphrase: optional[string], sealed: bytes) -> result[bytes, string]
    + returns the original plaintext
    - returns error when the private key cannot be loaded
    - returns error when the authentication tag does not match
    # unsealing
    -> std.ssh.parse_private_key
    -> std.crypto.rsa_decrypt_oaep
    -> std.crypto.aes_gcm_open
