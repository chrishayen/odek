# Requirement: "a package manager for machine-learning datasets and models"

Tracks named artifacts with versions, resolves dependencies, verifies content hashes, and maintains a local store. Remote fetching is injected by the caller.

std
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest
      # hashing
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + lowercase hex rendering
      # encoding
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on malformed input
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file at path
      - returns error when the file is missing
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating or replacing it
      # filesystem

mlpkg
  mlpkg.new_registry
    @ (store_path: string) -> registry_state
    + creates a registry backed by the given on-disk store path
    # construction
  mlpkg.load_index
    @ (state: registry_state) -> result[registry_state, string]
    + reads the persisted artifact index into memory
    - returns error when the index file is corrupt
    # persistence
    -> std.fs.read_all
    -> std.json.parse_object
  mlpkg.save_index
    @ (state: registry_state) -> result[void, string]
    + writes the in-memory index back to disk
    # persistence
    -> std.json.encode_object
    -> std.fs.write_all
  mlpkg.declare_artifact
    @ (state: registry_state, name: string, version: string, kind: string, deps: list[string]) -> registry_state
    + records metadata for a dataset or model version
    ? kind is "dataset" or "model"
    # metadata
  mlpkg.compute_digest
    @ (data: bytes) -> string
    + returns the hex SHA-256 digest for an artifact payload
    # hashing
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  mlpkg.install
    @ (state: registry_state, name: string, version: string, payload: bytes, expected_digest: string) -> result[registry_state, string]
    + stores the artifact under its digest and marks it installed
    - returns error when the computed digest does not match expected_digest
    # install
    -> std.crypto.sha256
    -> std.encoding.hex_encode
    -> std.fs.write_all
  mlpkg.resolve
    @ (state: registry_state, name: string, version: string) -> result[list[string], string]
    + returns a topological install order for the artifact and its dependencies
    - returns error on unknown names or dependency cycles
    # resolution
  mlpkg.locate
    @ (state: registry_state, name: string, version: string) -> optional[string]
    + returns the on-disk path for an installed artifact
    # lookup
  mlpkg.list_installed
    @ (state: registry_state) -> list[string]
    + returns every installed "name@version" pair
    # listing
  mlpkg.remove
    @ (state: registry_state, name: string, version: string) -> result[registry_state, string]
    + removes an artifact if nothing depends on it
    - returns error when other installed artifacts still reference it
    # removal
