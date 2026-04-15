# Requirement: "a git object store and history library"

A content-addressed object store for git-style blobs, trees, and commits, with ref management. Real git needs dozens of features; this exposes the core object graph operations.

std
  std.crypto
    std.crypto.sha1
      fn (data: bytes) -> bytes
      + returns the 20-byte SHA-1 digest
      # cryptography
  std.compression
    std.compression.deflate
      fn (data: bytes) -> bytes
      + zlib-deflates the input
      # compression
    std.compression.inflate
      fn (data: bytes) -> result[bytes, string]
      + zlib-inflates the input
      - returns error on corrupt stream
      # compression
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + returns lowercase hex representation
      # encoding
    std.encoding.hex_decode
      fn (s: string) -> result[bytes, string]
      + parses a lowercase or uppercase hex string
      - returns error on non-hex characters or odd length
      # encoding

git
  git.hash_blob
    fn (content: bytes) -> string
    + returns the git object id (hex sha1) for a blob with the standard "blob <len>\0" prefix
    # hashing
    -> std.crypto.sha1
    -> std.encoding.hex_encode
  git.write_blob
    fn (store: object_store, content: bytes) -> tuple[string, object_store]
    + stores a blob and returns its id plus the updated store
    # objects
    -> std.compression.deflate
  git.read_blob
    fn (store: object_store, id: string) -> result[bytes, string]
    + retrieves blob content by id
    - returns error when the id is unknown
    # objects
    -> std.compression.inflate
    -> std.encoding.hex_decode
  git.write_tree
    fn (store: object_store, entries: list[tree_entry]) -> tuple[string, object_store]
    + encodes a sorted tree of (mode, name, child_id) entries and stores it
    # objects
    -> std.crypto.sha1
    -> std.compression.deflate
  git.read_tree
    fn (store: object_store, id: string) -> result[list[tree_entry], string]
    + parses a stored tree object back to entries
    - returns error on malformed tree bytes
    # objects
    -> std.compression.inflate
  git.write_commit
    fn (store: object_store, tree_id: string, parents: list[string], author: string, message: string, timestamp: i64) -> tuple[string, object_store]
    + writes a commit object pointing at tree_id with the given parents
    # objects
    -> std.crypto.sha1
  git.read_commit
    fn (store: object_store, id: string) -> result[commit, string]
    + parses a commit object into tree, parents, author, message, timestamp
    - returns error when id does not refer to a commit
    # objects
  git.set_ref
    fn (store: object_store, name: string, target: string) -> object_store
    + points a named ref (e.g. "refs/heads/main") at the given object id
    # refs
  git.resolve_ref
    fn (store: object_store, name: string) -> optional[string]
    + returns the object id a ref points to, or none
    # refs
  git.walk_history
    fn (store: object_store, head: string) -> list[string]
    + returns commit ids reachable from head in first-parent order
    - returns empty list when head is unknown
    # traversal
  git.diff_trees
    fn (store: object_store, a: string, b: string) -> list[tree_change]
    + returns added, removed, and modified entries between two trees
    # diff
