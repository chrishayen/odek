# Requirement: "a library for creating, maintaining, finding, and reusing small modules across repositories"

A minimal component registry: modules are versioned units with content, metadata, and dependencies, searchable by name or tag.

std
  std.hash
    std.hash.sha256_hex
      fn (data: bytes) -> string
      + returns the hex-encoded SHA-256 digest
      # hashing
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns the current unix time in seconds
      # time

module_registry
  module_registry.new
    fn () -> registry_state
    + creates an empty in-memory registry
    # construction
  module_registry.publish
    fn (state: registry_state, name: string, content: bytes, tags: list[string], deps: list[string]) -> result[module_version, string]
    + stores a new version keyed by content hash and returns its metadata
    + assigns a monotonically increasing version number per name
    - returns error when name is empty
    # publication
    -> std.hash.sha256_hex
    -> std.time.now_seconds
  module_registry.get
    fn (state: registry_state, name: string, version: i32) -> optional[module_version]
    + returns the requested version when it exists
    - returns none when name or version is unknown
    # lookup
  module_registry.latest
    fn (state: registry_state, name: string) -> optional[module_version]
    + returns the highest version number for the given name
    # lookup
  module_registry.search_by_tag
    fn (state: registry_state, tag: string) -> list[module_version]
    + returns all latest versions whose tags contain the given tag
    # search
  module_registry.resolve_dependencies
    fn (state: registry_state, name: string, version: i32) -> result[list[module_version], string]
    + returns a topologically ordered list of transitive dependencies
    - returns error when a dependency cycle is detected
    - returns error when a dependency cannot be resolved
    # resolution
