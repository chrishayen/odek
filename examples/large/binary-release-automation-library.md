# Requirement: "a binary release automation library"

Automates building, packaging, and publishing compiled binaries across platforms. The project layer coordinates a pipeline; std provides generic build, archive, checksum, and upload primitives.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file fully into memory
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + creates or overwrites a file with the given bytes
      - returns error when the parent directory does not exist
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns entries in the directory, non-recursive
      - returns error when the path is not a directory
      # filesystem
  std.process
    std.process.run
      @ (program: string, args: list[string], env: map[string, string]) -> result[process_output, string]
      + runs a subprocess and captures stdout, stderr, exit code
      - returns error when the program cannot be launched
      # process
  std.archive
    std.archive.tar_gz
      @ (files: map[string, bytes]) -> bytes
      + produces a gzipped tar archive from a map of path to contents
      # archive
    std.archive.zip
      @ (files: map[string, bytes]) -> bytes
      + produces a zip archive from a map of path to contents
      # archive
  std.crypto
    std.crypto.sha256
      @ (data: bytes) -> bytes
      + returns the SHA-256 digest as 32 bytes
      # cryptography
  std.http
    std.http.post
      @ (url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an HTTP POST and returns status, headers, body
      - returns error on network failure
      # http
  std.encoding
    std.encoding.hex_encode
      @ (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding

releaser
  releaser.new_config
    @ (name: string, version: string, targets: list[target_spec]) -> release_config
    + creates a release configuration with project name, version, and target platforms
    # configuration
  releaser.build_target
    @ (cfg: release_config, target: target_spec) -> result[build_artifact, string]
    + invokes the platform-specific build command and captures the produced binary
    - returns error when the build command exits non-zero
    # build
    -> std.process.run
    -> std.fs.read_all
  releaser.package_artifact
    @ (artifact: build_artifact, format: string) -> result[bytes, string]
    + packages the artifact as tar.gz or zip based on format
    - returns error for unknown format
    # packaging
    -> std.archive.tar_gz
    -> std.archive.zip
  releaser.compute_checksums
    @ (packages: map[string, bytes]) -> map[string, string]
    + computes sha256 hex digest for each package keyed by filename
    # integrity
    -> std.crypto.sha256
    -> std.encoding.hex_encode
  releaser.write_manifest
    @ (cfg: release_config, checksums: map[string, string]) -> string
    + produces a manifest describing version, targets, and checksums
    # manifest
  releaser.upload_package
    @ (endpoint: string, auth_token: string, filename: string, data: bytes) -> result[void, string]
    + uploads a single package to a generic release endpoint with a bearer token
    - returns error when the endpoint rejects the upload
    # upload
    -> std.http.post
  releaser.run_pipeline
    @ (cfg: release_config, endpoint: string, auth_token: string) -> result[release_report, string]
    + runs build, package, checksum, and upload for every configured target
    - returns error on the first step that fails
    # orchestration
