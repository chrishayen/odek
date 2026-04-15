# Requirement: "a hashing and integrity tool for media files that records file hashes into a sidecar database and detects corruption or renames"

std
  std.fs
    std.fs.walk_dir
      fn (root: string) -> result[list[string], string]
      + returns absolute paths of all regular files beneath root
      - returns error when root is not readable
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the full contents of the file
      # filesystem
    std.fs.stat_mtime
      fn (path: string) -> result[i64, string]
      + returns last modification time in unix seconds
      # filesystem
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns the SHA-256 digest of data as 32 bytes
      # cryptography
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes to lowercase hexadecimal
      # encoding

media_hash
  media_hash.hash_file
    fn (path: string) -> result[string, string]
    + returns the lowercase hex SHA-256 of the file contents
    # hashing
    -> std.fs.read_all
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  media_hash.snapshot_tree
    fn (root: string) -> result[map[string, string], string]
    + walks root and returns a map from absolute path to file hash
    # snapshot
    -> std.fs.walk_dir
  media_hash.diff_snapshots
    fn (old: map[string, string], new: map[string, string]) -> snapshot_diff
    + classifies each path as unchanged, modified, added, or removed
    + detects renames by matching removed paths to added paths with equal hashes
    # integrity
  media_hash.verify_snapshot
    fn (snapshot: map[string, string]) -> list[string]
    + returns the list of paths whose on-disk content no longer matches the recorded hash
    # verification
  media_hash.serialize_snapshot
    fn (snapshot: map[string, string]) -> string
    + produces a line-oriented "hash  path" serialization
    # persistence
  media_hash.parse_snapshot
    fn (text: string) -> result[map[string, string], string]
    + parses a serialized snapshot back into a map
    - returns error on malformed lines
    # persistence
