# Requirement: "a password manager, filesystem compatible with pass"

Stores passwords as encrypted files in a directory tree. Entries are addressable by slash-separated names. Encryption and filesystem work happens in std primitives.

std
  std.fs
    std.fs.read_bytes
      fn (path: string) -> result[bytes, string]
      + returns the full contents of the file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_bytes
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the file, creating parent directories as needed
      - returns error when the path is not writable
      # filesystem
    std.fs.list_tree
      fn (root: string) -> result[list[string], string]
      + returns all file paths beneath root, relative to root
      - returns error when root does not exist
      # filesystem
    std.fs.remove
      fn (path: string) -> result[void, string]
      + removes the file at path
      - returns error when the file does not exist
      # filesystem
  std.crypto
    std.crypto.encrypt_for_recipients
      fn (plaintext: bytes, recipient_keys: list[bytes]) -> result[bytes, string]
      + returns a ciphertext blob encrypted to each recipient public key
      - returns error when the recipient list is empty
      # cryptography
    std.crypto.decrypt_with_identity
      fn (ciphertext: bytes, private_key: bytes) -> result[bytes, string]
      + returns the decrypted plaintext
      - returns error when the private key does not match any recipient
      - returns error when the ciphertext is truncated or tampered
      # cryptography
    std.crypto.random_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography

pass_store
  pass_store.open
    fn (root: string, recipients: list[bytes]) -> result[store_state, string]
    + returns a store_state bound to the given directory and recipient keys
    - returns error when root does not exist and cannot be created
    # construction
  pass_store.path_for
    fn (state: store_state, entry: string) -> string
    + returns the on-disk path for an entry name, appending the standard extension
    ? slashes in the entry name become subdirectories
    # addressing
  pass_store.insert
    fn (state: store_state, entry: string, secret: string) -> result[void, string]
    + encrypts the secret and writes it at the entry path
    - returns error when encryption fails
    # write
    -> std.crypto.encrypt_for_recipients
    -> std.fs.write_bytes
  pass_store.show
    fn (state: store_state, entry: string, private_key: bytes) -> result[string, string]
    + returns the decrypted secret for the entry
    - returns error when the entry does not exist
    - returns error when decryption fails
    # read
    -> std.fs.read_bytes
    -> std.crypto.decrypt_with_identity
  pass_store.list
    fn (state: store_state) -> result[list[string], string]
    + returns all entry names currently in the store
    # enumeration
    -> std.fs.list_tree
  pass_store.remove
    fn (state: store_state, entry: string) -> result[void, string]
    + removes the entry file from the store
    - returns error when the entry does not exist
    # delete
    -> std.fs.remove
  pass_store.generate
    fn (length: i32, alphabet: string) -> string
    + returns a random secret of the given length drawn from the alphabet
    ? length must be positive; alphabet must not be empty
    # generation
    -> std.crypto.random_bytes
