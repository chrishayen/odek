# Requirement: "a library for deploying code artifacts to a content delivery edge"

Uploads a directory tree to an object store, records a manifest of the release, and flips a pointer so edge caches begin serving the new version.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns every file path under root
      - returns error when root is not a directory
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the entire file contents as bytes
      # filesystem
  std.hash
    std.hash.sha256_hex
      fn (data: bytes) -> string
      + returns the SHA-256 digest of data as a lowercase hex string
      # hashing
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

deployer
  deployer.build_manifest
    fn (source_dir: string) -> result[manifest, string]
    + walks the source directory and returns a manifest listing every file with its hash and size
    - returns error when source_dir is missing
    # manifest
    -> std.fs.walk
    -> std.fs.read_all
    -> std.hash.sha256_hex
  deployer.upload_release
    fn (source_dir: string, release_id: string, store: object_store) -> result[manifest, string]
    + builds the manifest and uploads every file to the store under "releases/<release_id>/<path>"
    + skips upload when the store already has an object with the same hash
    - returns error when the store rejects an upload
    # distribution
    -> deployer.build_manifest
  deployer.write_pointer
    fn (release_id: string, store: object_store) -> result[void, string]
    + writes a "current" marker object naming the release_id
    - returns error when the store rejects the write
    # activation
    -> std.time.now_seconds
  deployer.deploy
    fn (source_dir: string, release_id: string, store: object_store) -> result[manifest, string]
    + uploads the release then writes the pointer
    - returns error from whichever step fails first
    # orchestration
    -> deployer.upload_release
    -> deployer.write_pointer
