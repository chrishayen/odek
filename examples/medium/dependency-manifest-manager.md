# Requirement: "a dependency manifest manager"

Loads a manifest file, adds and removes dependencies, resolves semver requirements against a registry, and writes a locked version file.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file's entire contents
      - returns error when the path does not exist
      # filesystem
    std.fs.write_atomic
      @ (path: string, data: bytes) -> result[void, string]
      + writes to a temp file and renames over the target
      # filesystem
  std.semver
    std.semver.parse
      @ (input: string) -> result[version, string]
      + parses a MAJOR.MINOR.PATCH string
      - returns error on malformed input
      # versioning
    std.semver.matches
      @ (requirement: string, v: version) -> bool
      + returns true when the version satisfies the requirement range (e.g. "^1.2", ">=0.4 <0.6")
      # versioning

dep_manager
  dep_manager.load
    @ (path: string) -> result[manifest_state, string]
    + parses a manifest file with direct dependencies
    - returns error on malformed manifest
    # manifest
    -> std.fs.read_all
  dep_manager.save
    @ (state: manifest_state, path: string) -> result[void, string]
    + serializes the manifest back to disk
    # manifest
    -> std.fs.write_atomic
  dep_manager.add
    @ (state: manifest_state, name: string, requirement: string) -> result[manifest_state, string]
    + adds or updates a dependency requirement
    - returns error when the requirement is not a valid semver range
    # edits
    -> std.semver.parse
  dep_manager.remove
    @ (state: manifest_state, name: string) -> result[manifest_state, string]
    + removes a dependency
    - returns error when the dependency is not in the manifest
    # edits
  dep_manager.resolve
    @ (state: manifest_state, registry_fn: string) -> result[map[string, string], string]
    + picks the highest version from a pluggable registry lookup that satisfies each requirement
    - returns error when no version satisfies a requirement
    # resolution
    -> std.semver.matches
  dep_manager.write_lock
    @ (resolved: map[string, string], path: string) -> result[void, string]
    + writes the resolved name-to-version map to a lockfile
    # lockfile
    -> std.fs.write_atomic
  dep_manager.check_outdated
    @ (state: manifest_state, registry_fn: string) -> list[tuple[string, string, string]]
    + returns (name, current, latest) for each dependency where the registry has a newer matching version
    # queries
