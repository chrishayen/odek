# Requirement: "a simplified directory-based authentication library"

A tiny user directory with hashed passwords and group membership, exposing authenticate and authorize calls.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest of data
      # cryptography
  std.rand
    std.rand.bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding

authdir
  authdir.new
    fn () -> directory_state
    + creates an empty directory with no users or groups
    # construction
  authdir.add_user
    fn (state: directory_state, username: string, password: string) -> result[directory_state, string]
    + stores the user with a fresh salt and a salted-sha256 hash
    - returns error when the username already exists
    # user_management
    -> std.rand.bytes
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  authdir.authenticate
    fn (state: directory_state, username: string, password: string) -> result[void, string]
    + returns ok when the password's salted hash matches the stored hash
    - returns error when the username is unknown
    - returns error when the password does not match
    # authentication
    -> std.crypto.sha256
  authdir.add_to_group
    fn (state: directory_state, username: string, group: string) -> result[directory_state, string]
    + adds the user to the given group
    - returns error when the username is unknown
    # authorization
  authdir.is_member
    fn (state: directory_state, username: string, group: string) -> bool
    + returns true when the user exists and belongs to the group
    # authorization
