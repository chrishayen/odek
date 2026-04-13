# Requirement: "a package manager for prebuilt binaries"

Manages a local store of installed binaries sourced from remote manifests.

std
  std.fs
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the path, creating it if needed
      - returns error on io failure
      # filesystem
    std.fs.remove
      @ (path: string) -> result[void, string]
      + removes the file at path
      - returns error when the path does not exist
      # filesystem
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the hex-encoded sha256 digest
      # hashing
  std.net
    std.net.http_get
      @ (url: string) -> result[bytes, string]
      + performs an http get and returns the response body
      - returns error on non-2xx status
      # network

pkgmgr
  pkgmgr.new_store
    @ (root_dir: string) -> store_state
    + creates a store rooted at the given directory
    # construction
  pkgmgr.parse_manifest
    @ (raw: string) -> result[manifest, string]
    + parses a manifest describing name, version, url, and sha256
    - returns error when required fields are missing
    # manifest
  pkgmgr.install
    @ (s: store_state, m: manifest) -> result[store_state, string]
    + downloads the binary, verifies its digest, and records it in the store
    - returns error when the downloaded digest does not match
    # install
    -> std.net.http_get
    -> std.hash.sha256_hex
    -> std.fs.write_all
  pkgmgr.uninstall
    @ (s: store_state, name: string) -> result[store_state, string]
    + removes the binary and its record
    - returns error when the package is not installed
    # uninstall
    -> std.fs.remove
  pkgmgr.list_installed
    @ (s: store_state) -> list[installed_entry]
    + returns name, version, and install path for each installed package
    # query
  pkgmgr.find_installed
    @ (s: store_state, name: string) -> optional[installed_entry]
    + returns the installed entry for the given name
    - returns none when not installed
    # query
