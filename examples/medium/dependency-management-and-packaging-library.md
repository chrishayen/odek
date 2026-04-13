# Requirement: "a dependency management and packaging library"

Models a project manifest, resolves a dependency graph against a pluggable registry, and produces a lockfile. Registry IO is a std primitive.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, overwriting
      - returns error on permission failure
      # filesystem
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.semver
    std.semver.parse
      @ (raw: string) -> result[semver, string]
      + parses a "major.minor.patch" version
      - returns error on non-semver input
      # versioning
    std.semver.satisfies
      @ (v: semver, constraint: string) -> bool
      + returns true if v meets a range expression such as ">=1.2.0 <2.0.0"
      # versioning

deps
  deps.load_manifest
    @ (path: string) -> result[manifest, string]
    + reads a manifest file and parses its name, version, and dependency constraints
    - returns error when the file is missing or malformed
    # manifest
    -> std.fs.read_all
    -> std.json.parse_object
    -> std.semver.parse
  deps.add_dependency
    @ (m: manifest, name: string, constraint: string) -> result[manifest, string]
    + records a new dependency constraint on the manifest
    - returns error when constraint is not a valid range
    # manifest
  deps.remove_dependency
    @ (m: manifest, name: string) -> manifest
    + removes a dependency if present
    # manifest
  deps.resolve
    @ (m: manifest, registry: map[string, list[string]]) -> result[map[string, string], string]
    + returns a name-to-version map where every dependency's chosen version satisfies its constraint
    - returns error when no compatible set exists
    # resolution
    -> std.semver.satisfies
  deps.write_lockfile
    @ (path: string, resolved: map[string, string]) -> result[void, string]
    + serializes resolved name-version pairs to a lockfile
    - returns error on write failure
    # lockfile
    -> std.json.encode_object
    -> std.fs.write_all
  deps.read_lockfile
    @ (path: string) -> result[map[string, string], string]
    + loads a lockfile and returns its name-version map
    - returns error when the file is missing or malformed
    # lockfile
    -> std.fs.read_all
    -> std.json.parse_object
