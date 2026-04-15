# Requirement: "a combined virtual environment and package manager"

Maintains isolated environments with declared dependencies. Package resolution and filesystem layout work happens in std primitives.

std
  std.fs
    std.fs.mkdir_all
      fn (path: string) -> result[void, string]
      + creates a directory and any missing parents
      - returns error when the parent is a file
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, overwriting
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file to bytes
      - returns error when the file is missing
      # filesystem
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string map as a JSON object
      # serialization
  std.hash
    std.hash.sha256_hex
      fn (data: bytes) -> string
      + returns the SHA-256 hex digest
      # hashing

envmgr
  envmgr.create
    fn (root: string, name: string) -> result[env_handle, string]
    + creates a new environment directory with an empty manifest
    - returns error when an environment with that name already exists
    # construction
    -> std.fs.mkdir_all
    -> std.json.encode_object
    -> std.fs.write_all
  envmgr.load
    fn (root: string, name: string) -> result[env_handle, string]
    + loads an existing environment from disk
    - returns error when the environment does not exist
    # loading
    -> std.fs.read_all
    -> std.json.parse_object
  envmgr.add_package
    fn (env: env_handle, name: string, version: string) -> result[env_handle, string]
    + records a dependency in the manifest
    - returns error when the version string is empty
    # dependencies
  envmgr.remove_package
    fn (env: env_handle, name: string) -> env_handle
    + removes a package from the manifest if present
    # dependencies
  envmgr.resolve
    fn (env: env_handle, index: map[string, list[string]]) -> result[map[string, string], string]
    + picks a concrete version for each dependency given an index of available versions
    - returns error when a declared dependency is not in the index
    # resolution
  envmgr.install
    fn (env: env_handle, resolved: map[string, string]) -> result[env_handle, string]
    + materializes resolved packages into the environment directory with content hashing
    - returns error when a package file cannot be written
    # installation
    -> std.hash.sha256_hex
    -> std.fs.mkdir_all
    -> std.fs.write_all
  envmgr.save
    fn (env: env_handle) -> result[void, string]
    + writes the manifest back to disk
    # persistence
    -> std.json.encode_object
    -> std.fs.write_all
