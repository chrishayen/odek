# Requirement: "a library to obfuscate source scripts, bind them to a specific machine, and expire them"

Transforms source text into an obfuscated loader whose payload is encrypted. The loader checks a machine fingerprint and an expiry timestamp before executing.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when the file cannot be opened
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes the contents to the given path, replacing any existing file
      # filesystem
  std.crypto
    std.crypto.aes_gcm_encrypt
      fn (key: bytes, nonce: bytes, plaintext: bytes) -> result[bytes, string]
      + returns ciphertext with authentication tag appended
      - returns error when key length is invalid
      # cryptography
    std.crypto.aes_gcm_decrypt
      fn (key: bytes, nonce: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns plaintext when the tag verifies
      - returns error when the tag does not verify
      # cryptography
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns the SHA-256 digest of data
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + returns standard base64 text
      # encoding
    std.encoding.base64_decode
      fn (encoded: string) -> result[bytes, string]
      + returns decoded bytes
      - returns error on invalid input
      # encoding
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.host
    std.host.machine_fingerprint
      fn () -> string
      + returns a stable identifier derived from host properties
      # host

script_protect
  script_protect.derive_key
    fn (password: string, salt: bytes) -> bytes
    + returns a 32-byte key derived from the password and salt
    # crypto
    -> std.crypto.sha256
  script_protect.compute_binding_hash
    fn (fingerprint: string, salt: bytes) -> bytes
    + returns a digest that uniquely identifies the allowed host for this artifact
    # binding
    -> std.crypto.sha256
  script_protect.encrypt_payload
    fn (source: string, key: bytes) -> result[encrypted_payload, string]
    + returns an encrypted payload containing nonce and ciphertext
    - returns error when encryption fails
    # crypto
    -> std.crypto.random_bytes
    -> std.crypto.aes_gcm_encrypt
  script_protect.decrypt_payload
    fn (payload: encrypted_payload, key: bytes) -> result[string, string]
    + returns the original source when the key and tag are correct
    - returns error when decryption fails
    # crypto
    -> std.crypto.aes_gcm_decrypt
  script_protect.build_artifact
    fn (source: string, password: string, bound_fingerprint: optional[string], expires_at: optional[i64]) -> result[string, string]
    + returns a self-contained loader source whose payload is encrypted and whose header carries the optional binding hash and expiry
    - returns error when encryption fails
    # packaging
    -> std.crypto.random_bytes
    -> std.encoding.base64_encode
  script_protect.parse_artifact
    fn (artifact: string) -> result[artifact_header, string]
    + returns the header describing payload, binding hash, and expiry
    - returns error when the artifact is not recognizably produced by build_artifact
    # packaging
    -> std.encoding.base64_decode
  script_protect.check_binding
    fn (header: artifact_header) -> result[void, string]
    + returns ok when the artifact has no binding or when the current host matches the bound fingerprint
    - returns error when the fingerprints do not match
    # policy
    -> std.host.machine_fingerprint
  script_protect.check_expiry
    fn (header: artifact_header) -> result[void, string]
    + returns ok when no expiry is set or the current time is before expiry
    - returns error when the artifact has expired
    # policy
    -> std.time.now_seconds
  script_protect.load_artifact
    fn (artifact: string, password: string) -> result[string, string]
    + parses the header, checks binding and expiry, then decrypts and returns the original source
    - returns error on any policy failure or decryption failure
    # loading
  script_protect.write_artifact
    fn (source_path: string, output_path: string, password: string, bind: bool, ttl_seconds: optional[i64]) -> result[void, string]
    + reads the source, builds a protected artifact, and writes it to the output path
    - returns error when reading, building, or writing fails
    # orchestration
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.time.now_seconds
