# Requirement: "a system package manager library"

Resolves dependencies, tracks installed packages, and orchestrates install and remove operations. The actual file staging is delegated through a std filesystem primitive.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the full contents of a regular file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to a file
      # filesystem
    std.fs.remove
      fn (path: string) -> result[void, string]
      + removes a file
      # filesystem
  std.encoding
    std.encoding.sha256_hex
      fn (data: bytes) -> string
      + returns the hex-encoded SHA-256 of data
      # hashing
  std.archive
    std.archive.list_entries
      fn (archive: bytes) -> result[list[archive_entry], string]
      + returns the names and sizes of entries in an archive
      - returns error on malformed archive
      # archive
    std.archive.extract_entry
      fn (archive: bytes, entry_name: string) -> result[bytes, string]
      + returns the bytes of a single entry
      # archive

pkg
  pkg.parse_manifest
    fn (raw: string) -> result[manifest, string]
    + parses a manifest listing name, version, and dependencies
    - returns error when name or version field is missing
    # manifest
  pkg.new_index
    fn () -> index_state
    + creates an empty installed-package index
    # construction
  pkg.index_add
    fn (index: index_state, name: string, version: string, file_list: list[string]) -> index_state
    + records an installed package and the files it owns
    # registration
  pkg.index_remove
    fn (index: index_state, name: string) -> index_state
    + removes a package from the index
    - leaves the index unchanged when name is not installed
    # registration
  pkg.resolve_deps
    fn (want: manifest, available: map[string, manifest]) -> result[list[string], string]
    + returns a topologically ordered install list
    - returns error on unsatisfiable dependency
    - returns error on dependency cycle
    # resolution
  pkg.verify_archive
    fn (archive: bytes, expected_hash: string) -> result[void, string]
    + returns ok when the sha256 of archive matches expected_hash
    - returns error on hash mismatch
    # integrity
    -> std.encoding.sha256_hex
  pkg.install
    fn (index: index_state, manifest: manifest, archive: bytes, install_root: string) -> result[index_state, string]
    + extracts archive entries under install_root and records them in the index
    - returns error when extraction fails
    - returns error when a file conflicts with an already-installed package
    # install
    -> std.archive.list_entries
    -> std.archive.extract_entry
    -> std.fs.write_all
  pkg.remove
    fn (index: index_state, name: string) -> result[index_state, string]
    + deletes files owned by name and removes it from the index
    - returns error when another installed package depends on name
    # remove
    -> std.fs.remove
  pkg.list_installed
    fn (index: index_state) -> list[tuple[string, string]]
    + returns name and version pairs for all installed packages
    # query
  pkg.find_owner
    fn (index: index_state, path: string) -> optional[string]
    + returns the package that owns the given installed path
    # query
