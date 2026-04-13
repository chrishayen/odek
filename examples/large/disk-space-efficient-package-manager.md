# Requirement: "a disk space efficient package manager"

Stores each package version once in a content-addressed store and materializes project dependency trees via hard links instead of copies.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file fully into memory
      - returns error when path is missing
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating parent directories as needed
      # filesystem
    std.fs.hardlink
      @ (target: string, link: string) -> result[void, string]
      + creates a hard link at link pointing to target
      - returns error when the target and link are on different filesystems
      # filesystem
    std.fs.mkdir_all
      @ (path: string) -> result[void, string]
      + creates path and any missing parents
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns directory entries
      # filesystem
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the sha-256 digest
      # cryptography
  std.net
    std.net.http_get
      @ (url: string) -> result[bytes, string]
      + fetches the body at url
      - returns error on non-2xx or transport failure
      # networking
  std.archive
    std.archive.extract_tar_gz
      @ (data: bytes) -> result[list[archive_entry], string]
      + decompresses and extracts a gzipped tar into in-memory entries
      - returns error on corrupt archive
      # archives
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a json object into a string map
      - returns error on invalid json
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + serializes a string map as a json object
      # serialization

pkg_manager
  pkg_manager.new
    @ (store_root: string, project_root: string) -> pkg_manager_state
    + returns a manager rooted at the given store and project directories
    # construction
  pkg_manager.parse_manifest
    @ (raw: string) -> result[manifest, string]
    + parses a project manifest listing dependencies with name and version
    - returns error when required fields are missing
    # manifest
    -> std.json.parse_object
  pkg_manager.resolve_dependencies
    @ (state: pkg_manager_state, manifest: manifest) -> result[list[resolved_pkg], string]
    + walks each dependency and returns the full flattened closure with versions pinned
    - returns error on version conflicts that have no valid resolution
    ? transitive dependencies are read from fetched package metadata
    # resolution
  pkg_manager.fetch_package
    @ (state: pkg_manager_state, name: string, version: string) -> result[bytes, string]
    + downloads the archive for name@version from the registry
    - returns error on http failure
    # fetching
    -> std.net.http_get
  pkg_manager.store_package
    @ (state: pkg_manager_state, archive: bytes) -> result[string, string]
    + extracts the archive, addresses each file by its sha-256 content hash
    + writes each unique file once into the content store and returns the package content id
    ? files that already exist in the store are not rewritten
    # content_store
    -> std.archive.extract_tar_gz
    -> std.crypto.sha256
    -> std.fs.write_all
  pkg_manager.materialize_project
    @ (state: pkg_manager_state, resolved: list[resolved_pkg]) -> result[void, string]
    + creates the project's dependency directory tree using hard links into the content store
    - returns error when the store and project are on different filesystems
    # materialization
    -> std.fs.hardlink
    -> std.fs.mkdir_all
  pkg_manager.install
    @ (state: pkg_manager_state, manifest_path: string) -> result[void, string]
    + full install: parse manifest, resolve, fetch missing, store, materialize
    # orchestration
    -> std.fs.read_all
  pkg_manager.prune_store
    @ (state: pkg_manager_state, referenced_ids: list[string]) -> result[i64, string]
    + removes content store entries not in referenced_ids and returns bytes freed
    # gc
    -> std.fs.list_dir
