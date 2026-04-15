# Requirement: "a duplicate file finder and remover"

Walks a directory, groups files by size then by content hash, and reports or deletes duplicates.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns paths of every regular file beneath root
      - returns error when root does not exist
      # fs
    std.fs.stat_size
      fn (path: string) -> result[i64, string]
      + returns the byte size of the file
      - returns error when path is not a regular file
      # fs
    std.fs.remove
      fn (path: string) -> result[void, string]
      + deletes the file at path
      - returns error on permission failure
      # fs
    std.fs.open_read
      fn (path: string) -> result[file_handle, string]
      + opens a file for streaming reads
      - returns error when file cannot be opened
      # fs
    std.fs.read_chunk
      fn (handle: file_handle, size: i32) -> result[bytes, string]
      + reads up to size bytes from the handle
      - returns error on I/O failure
      # fs
  std.crypto
    std.crypto.blake3_new
      fn () -> hasher_state
      + creates a streaming BLAKE3 hasher
      # cryptography
    std.crypto.blake3_update
      fn (hasher: hasher_state, data: bytes) -> hasher_state
      + feeds data into the hasher
      # cryptography
    std.crypto.blake3_finalize
      fn (hasher: hasher_state) -> bytes
      + returns the final 32-byte digest
      # cryptography

dedup
  dedup.scan
    fn (root: string) -> result[list[string], string]
    + returns every regular file under root
    - returns error when the root cannot be walked
    # discovery
    -> std.fs.walk
  dedup.group_by_size
    fn (paths: list[string]) -> result[map[i64, list[string]], string]
    + groups paths by file size
    - returns error when a path cannot be stat'd
    # grouping
    -> std.fs.stat_size
  dedup.quick_fingerprint
    fn (path: string) -> result[bytes, string]
    + hashes the first and last page of the file as a cheap prefilter
    - returns error on I/O failure
    # hashing
    -> std.fs.open_read
    -> std.fs.read_chunk
    -> std.crypto.blake3_new
    -> std.crypto.blake3_update
    -> std.crypto.blake3_finalize
  dedup.full_hash
    fn (path: string) -> result[bytes, string]
    + hashes the entire file contents
    - returns error on I/O failure
    # hashing
    -> std.fs.open_read
    -> std.fs.read_chunk
    -> std.crypto.blake3_new
    -> std.crypto.blake3_update
    -> std.crypto.blake3_finalize
  dedup.find_duplicates
    fn (root: string) -> result[list[list[string]], string]
    + returns groups of paths that share identical content
    - returns error when scanning fails
    # pipeline
  dedup.pick_keeper
    fn (group: list[string]) -> string
    + selects one path from a duplicate group to retain
    ? chooses the shortest path; ties broken lexicographically
    # policy
  dedup.remove_duplicates
    fn (groups: list[list[string]]) -> result[i32, string]
    + deletes all non-keeper files and returns the count removed
    - returns error when any deletion fails
    # removal
    -> std.fs.remove
  dedup.report
    fn (groups: list[list[string]]) -> string
    + returns a human-readable summary of duplicate groups
    # reporting
