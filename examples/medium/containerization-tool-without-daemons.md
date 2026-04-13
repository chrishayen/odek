# Requirement: "a daemonless container builder and runner"

Builds an isolated process environment from a declarative spec and runs commands inside it, without a long-running daemon. Filesystem and process isolation are abstracted behind handles.

std
  std.fs
    std.fs.ensure_dir
      @ (path: string) -> result[void, string]
      + creates the directory and any missing parents
      # filesystem
    std.fs.copy_tree
      @ (src: string, dst: string) -> result[void, string]
      + copies a directory tree recursively
      - returns error when src does not exist
      # filesystem
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex SHA-256 digest
      # hashing

sandbox
  sandbox.parse_spec
    @ (source: string) -> result[container_spec, string]
    + parses a declarative spec of base image, mounts, env vars, and command
    - returns error on malformed input
    # parsing
  sandbox.resolve_spec
    @ (spec: container_spec, cache_dir: string) -> result[resolved_spec, string]
    + materializes a deterministic root directory for the spec in the cache
    - returns error when a referenced mount source is missing
    # preparation
    -> std.fs.ensure_dir
    -> std.fs.copy_tree
    -> std.hash.sha256_hex
  sandbox.run
    @ (resolved: resolved_spec, command: list[string], isolation: isolation_handle) -> result[i32, string]
    + executes the command inside the isolated environment and returns the exit code
    - returns error when the command cannot be launched
    # execution
  sandbox.cleanup
    @ (resolved: resolved_spec, keep_cache: bool) -> result[void, string]
    + removes per-run state; optionally keeps the cached root directory
    # cleanup
  sandbox.list_cached
    @ (cache_dir: string) -> result[list[string], string]
    + returns the ids of cached root directories
    # cache
