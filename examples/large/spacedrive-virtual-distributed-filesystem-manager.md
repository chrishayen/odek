# Requirement: "a file manager built on a virtual distributed filesystem"

Files are addressed by content hash, replicated across locations, and browsed through a virtual directory tree stored independently from storage locations.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + computes SHA-256 digest
      # cryptography
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads an entire file
      # io
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file
      # io
    std.fs.list_dir
      fn (dir: string) -> result[list[string], string]
      + returns entry names inside dir
      # io
  std.kv
    std.kv.put
      fn (key: string, value: bytes) -> result[void, string]
      + writes a key-value pair
      # storage
    std.kv.get
      fn (key: string) -> result[optional[bytes], string]
      + returns the value or none
      # storage
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

spacedrive
  spacedrive.hash_file
    fn (path: string) -> result[string, string]
    + returns the hex-encoded SHA-256 of the file contents
    # content_addressing
    -> std.fs.read_all
    -> std.crypto.sha256
  spacedrive.register_location
    fn (name: string, root: string) -> result[location_id, string]
    + registers a storage location rooted at a local path
    - returns error when the root is not readable
    # locations
    -> std.fs.list_dir
    -> std.kv.put
  spacedrive.index_location
    fn (location: location_id) -> result[i32, string]
    + walks the location, hashes new files, and records them in the catalog
    + returns the number of files indexed
    # indexing
    -> std.fs.list_dir
    -> std.kv.put
  spacedrive.catalog_entry
    fn (file_hash: string, name: string, size: i64) -> result[void, string]
    + stores a catalog entry keyed by content hash
    # catalog
    -> std.time.now_seconds
    -> std.kv.put
  spacedrive.find_by_hash
    fn (file_hash: string) -> result[optional[file_entry], string]
    + returns the catalog entry for a content hash
    # catalog
    -> std.kv.get
  spacedrive.list_replicas
    fn (file_hash: string) -> result[list[location_id], string]
    + returns every location known to store the hash
    # replication
    -> std.kv.get
  spacedrive.replicate
    fn (file_hash: string, target: location_id) -> result[void, string]
    + copies the file to the target location and records the replica
    - returns error when no source replica can be read
    # replication
    -> std.fs.read_all
    -> std.fs.write_all
    -> std.kv.put
  spacedrive.create_virtual_dir
    fn (parent: virtual_path, name: string) -> result[virtual_path, string]
    + creates a virtual directory node under parent
    - returns error when parent does not exist
    # virtual_tree
    -> std.kv.put
  spacedrive.link_into_virtual
    fn (dir: virtual_path, name: string, file_hash: string) -> result[void, string]
    + creates a named reference to a cataloged file inside a virtual directory
    - returns error when the file hash is unknown
    # virtual_tree
    -> std.kv.get
    -> std.kv.put
  spacedrive.list_virtual
    fn (dir: virtual_path) -> result[list[virtual_entry], string]
    + returns entries inside a virtual directory
    # virtual_tree
    -> std.kv.get
  spacedrive.resolve_virtual
    fn (path: virtual_path) -> result[file_entry, string]
    + resolves a virtual path to its backing file entry
    - returns error when the path is a directory or does not exist
    # virtual_tree
