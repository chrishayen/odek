# Requirement: "a toolchain version manager that tracks project dependencies"

Installs multiple toolchain versions side by side, records a project's pinned version, and lists its direct dependencies.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the complete contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes the bytes, creating parent directories as needed
      - returns error when the destination is not writable
      # filesystem
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

toolchain
  toolchain.install_version
    fn (root: string, version: string, archive: bytes) -> result[void, string]
    + extracts the archive under root/version and marks the install complete
    - returns error when archive is empty
    # installation
    -> std.fs.write_all
  toolchain.list_installed
    fn (root: string) -> result[list[string], string]
    + returns the versions present under root
    - returns error when root does not exist
    # installation
  toolchain.pin_project
    fn (project_dir: string, version: string) -> result[void, string]
    + writes the pinned version into the project's version file
    # project_pinning
    -> std.fs.write_all
  toolchain.active_version
    fn (project_dir: string) -> result[string, string]
    + returns the version pinned for the project
    - returns error when no pin file is present
    # project_pinning
    -> std.fs.read_all
  toolchain.read_manifest
    fn (project_dir: string) -> result[map[string, string], string]
    + returns the project's direct dependencies as a name-to-version map
    - returns error when the manifest is missing or invalid
    # dependencies
    -> std.fs.read_all
    -> std.json.parse_object
  toolchain.write_manifest
    fn (project_dir: string, deps: map[string, string]) -> result[void, string]
    + writes the dependency map back to the manifest file
    # dependencies
    -> std.json.encode_object
    -> std.fs.write_all
