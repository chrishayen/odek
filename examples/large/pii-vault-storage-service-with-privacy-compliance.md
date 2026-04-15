# Requirement: "a PII storage service with privacy compliance features"

A vault for personal data with encryption at rest, consent tracking, audit logging, and right-to-erasure support.

std
  std.crypto
    std.crypto.aes_gcm_encrypt
      fn (key: bytes, plaintext: bytes) -> result[bytes, string]
      + returns ciphertext prefixed with a random nonce
      - returns error when key is not 32 bytes
      # cryptography
    std.crypto.aes_gcm_decrypt
      fn (key: bytes, ciphertext: bytes) -> result[bytes, string]
      + returns plaintext when authentication tag validates
      - returns error when ciphertext is tampered
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
  std.hash
    std.hash.sha256
      fn (data: bytes) -> bytes
      + returns 32-byte digest
      # hashing
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time
  std.uuid
    std.uuid.v4
      fn () -> string
      + returns a random UUID v4 string
      # identifiers

pii_vault
  pii_vault.open
    fn (master_key: bytes) -> result[vault_state, string]
    + creates a vault instance bound to the given master key
    - returns error when master_key is not 32 bytes
    # construction
  pii_vault.store_record
    fn (state: vault_state, subject_id: string, fields: map[string,string], consent_tags: list[string]) -> result[string, string]
    + encrypts each field and returns a record id
    + records a lookup token so fields can be searched without decryption
    # storage
    -> std.crypto.aes_gcm_encrypt
    -> std.hash.sha256
    -> std.uuid.v4
    -> std.time.now_seconds
  pii_vault.fetch_record
    fn (state: vault_state, record_id: string) -> result[map[string,string], string]
    + decrypts and returns all fields for a record
    - returns error when record_id is unknown
    # retrieval
    -> std.crypto.aes_gcm_decrypt
  pii_vault.find_by_field
    fn (state: vault_state, field_name: string, field_value: string) -> list[string]
    + returns record ids whose deterministic lookup token matches
    # search
    -> std.hash.sha256
  pii_vault.update_consent
    fn (state: vault_state, record_id: string, consent_tags: list[string]) -> result[void, string]
    + replaces the consent tags on a stored record
    - returns error when record_id is unknown
    # consent
    -> std.time.now_seconds
  pii_vault.forget_subject
    fn (state: vault_state, subject_id: string) -> i32
    + zeroes out all records for subject_id and returns the count removed
    # erasure
    -> std.time.now_seconds
  pii_vault.export_subject
    fn (state: vault_state, subject_id: string) -> result[list[map[string,string]], string]
    + returns all decrypted records belonging to subject_id
    # portability
    -> std.crypto.aes_gcm_decrypt
  pii_vault.audit_log
    fn (state: vault_state, since_seconds: i64) -> list[audit_entry]
    + returns all audit entries since the given unix time
    # auditing
