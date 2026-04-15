# Requirement: "a library for building, deploying, and managing application packages"

A package lifecycle library: build a source tree into an immutable artifact, push it to a remote store, install it on a host, and run lifecycle hooks (init, start, stop).

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the entire file contents as bytes
      - returns error when the path is missing
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to the file, creating or overwriting
      # filesystem
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + returns every file path under root
      - returns error when root is not a directory
      # filesystem
  std.archive
    std.archive.pack
      fn (files: map[string, bytes]) -> bytes
      + returns a single byte blob containing every file keyed by relative path
      # archive
    std.archive.unpack
      fn (blob: bytes) -> result[map[string, bytes], string]
      + returns the files encoded in the blob
      - returns error on a truncated or malformed blob
      # archive
  std.hash
    std.hash.sha256_hex
      fn (data: bytes) -> string
      + returns the SHA-256 digest of data as a lowercase hex string
      # hashing
  std.os
    std.os.run_command
      fn (program: string, args: list[string], env: map[string, string]) -> result[process_result, string]
      + runs the program and returns exit code, stdout, and stderr
      # process

package_manager
  package_manager.build_artifact
    fn (source_dir: string, name: string, version: string) -> result[artifact, string]
    + walks the source directory, archives every file, and returns an artifact with name, version, and content hash
    - returns error when source_dir is empty or missing
    # build
    -> std.fs.walk
    -> std.fs.read_all
    -> std.archive.pack
    -> std.hash.sha256_hex
  package_manager.write_artifact
    fn (art: artifact, dest_path: string) -> result[void, string]
    + writes the artifact blob to dest_path
    # build
    -> std.fs.write_all
  package_manager.publish
    fn (art: artifact, store: artifact_store) -> result[void, string]
    + uploads the artifact to the pluggable store under "<name>/<version>"
    - returns error when an artifact with the same name and version already exists
    # distribution
  package_manager.fetch
    fn (name: string, version: string, store: artifact_store) -> result[artifact, string]
    + downloads an artifact from the store
    - returns error when the store has no matching artifact
    # distribution
  package_manager.install
    fn (art: artifact, install_dir: string) -> result[void, string]
    + unpacks the artifact into install_dir, preserving relative paths
    - returns error when install_dir already contains files for this name+version
    # installation
    -> std.archive.unpack
    -> std.fs.write_all
  package_manager.run_hook
    fn (install_dir: string, hook: string) -> result[void, string]
    + runs the named lifecycle script ("init", "start", "stop", "health") from install_dir
    - returns error when the hook script exits non-zero
    - returns error when the hook script does not exist
    # lifecycle
    -> std.os.run_command
  package_manager.start
    fn (install_dir: string) -> result[void, string]
    + runs the init hook then the start hook
    # lifecycle
    -> package_manager.run_hook
  package_manager.stop
    fn (install_dir: string) -> result[void, string]
    + runs the stop hook
    # lifecycle
    -> package_manager.run_hook
