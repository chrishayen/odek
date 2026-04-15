# Requirement: "a library for fetching, verifying, and installing signed artifacts from a trusted repository"

Downloads an artifact, verifies its signature against a trusted public key, and swaps it into place atomically.

std
  std.http
    std.http.get
      fn (url: string) -> result[bytes, string]
      + returns the response body for a successful GET
      - returns error on non-2xx status or transport failure
      # http
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns the 32-byte SHA-256 digest of data
      # cryptography
    std.crypto.ed25519_verify
      fn (public_key: bytes, message: bytes, signature: bytes) -> bool
      + returns true when the signature is valid for the message under the key
      - returns false on any tampering or wrong key
      # cryptography
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the full file contents
      - returns error when the path does not exist
      # filesystem
    std.fs.atomic_write
      fn (path: string, data: bytes) -> result[void, string]
      + writes to a temporary file and renames it over the target
      - returns error when the directory is not writable
      # filesystem

signed_artifacts
  signed_artifacts.fetch_manifest
    fn (repo_url: string, name: string) -> result[manifest, string]
    + downloads and parses the manifest describing artifact url, digest, and signature
    - returns error when the repo returns no manifest for the name
    # manifest
    -> std.http.get
  signed_artifacts.verify_artifact
    fn (data: bytes, manifest: manifest, trusted_key: bytes) -> result[void, string]
    + returns ok when the digest matches and the signature verifies under the trusted key
    - returns error when the digest does not match the manifest
    - returns error when the signature is not valid
    # verification
    -> std.crypto.sha256
    -> std.crypto.ed25519_verify
  signed_artifacts.install
    fn (repo_url: string, name: string, install_path: string, trusted_key: bytes) -> result[void, string]
    + fetches, verifies, and atomically writes the artifact to the install path
    - returns error when verification fails, leaving the existing file untouched
    # install
    -> std.http.get
    -> std.fs.atomic_write
  signed_artifacts.current_version
    fn (install_path: string) -> result[string, string]
    + returns the version string embedded in the installed artifact header
    - returns error when the path does not contain a recognized artifact
    # inspection
    -> std.fs.read_all
