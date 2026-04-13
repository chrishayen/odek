# Requirement: "a library to synchronize installed packages across multiple machines"

Tracks package installations locally, reconciles with a shared manifest, and reports what to install or remove on each machine. No network or package manager is actually invoked; those are the caller's job.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of the file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the path
      # filesystem
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a map as JSON
      # serialization

pkgsync
  pkgsync.load_manifest
    @ (path: string) -> result[manifest_state, string]
    + reads the shared manifest file
    - returns error when the file is missing or malformed
    # storage
    -> std.fs.read_all
    -> std.json.parse_object
  pkgsync.save_manifest
    @ (manifest: manifest_state, path: string) -> result[void, string]
    + writes the manifest atomically
    # storage
    -> std.json.encode_object
    -> std.fs.write_all
  pkgsync.record_install
    @ (manifest: manifest_state, machine: string, package: string) -> manifest_state
    + adds the package under the machine's installed set
    # mutation
  pkgsync.record_remove
    @ (manifest: manifest_state, machine: string, package: string) -> manifest_state
    + removes the package from the machine's installed set
    # mutation
  pkgsync.diff
    @ (manifest: manifest_state, machine: string, local_installed: list[string]) -> sync_plan
    + returns the set of packages to install and to remove on the given machine
    ? local_installed is what the caller found actually present on the machine
    # reconciliation
  pkgsync.merge_peer
    @ (manifest: manifest_state, peer: manifest_state) -> manifest_state
    + unions entries from a peer manifest into the local one
    # reconciliation
  pkgsync.list_packages_for_machine
    @ (manifest: manifest_state, machine: string) -> list[string]
    + returns the recorded package list for the machine
    # query
