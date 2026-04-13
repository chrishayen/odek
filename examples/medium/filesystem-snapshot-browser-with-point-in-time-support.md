# Requirement: "a file-level snapshot browser for filesystems that support point-in-time snapshots"

Lets a user list, diff, and restore prior versions of a file from filesystem snapshots. The library does not know which snapshotting technology is underneath; callers provide a snapshot root lister.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file at path
      - returns error when path does not exist or is unreadable
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating or replacing
      - returns error on permission or io failure
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns the names of entries directly inside path
      - returns error when path is not a directory
      # filesystem
    std.fs.stat_mtime
      @ (path: string) -> result[i64, string]
      + returns modification time as unix ms
      - returns error when path does not exist
      # filesystem

snapshot_browser
  snapshot_browser.discover_versions
    @ (live_path: string, snapshot_roots: list[string]) -> result[list[file_version], string]
    + returns one version per snapshot root that contains a copy of live_path
    ? each snapshot root mirrors the live tree under the same relative path
    # discovery
    -> std.fs.list_dir
    -> std.fs.stat_mtime
  snapshot_browser.load_version
    @ (version: file_version) -> result[bytes, string]
    + returns the bytes of the file at that version
    - returns error when the version file is missing
    # loading
    -> std.fs.read_all
  snapshot_browser.diff_versions
    @ (a: bytes, b: bytes) -> list[diff_hunk]
    + returns line-level hunks between two byte sequences
    + returns an empty list when the inputs are identical
    # diffing
  snapshot_browser.restore_version
    @ (version: file_version, target_path: string) -> result[void, string]
    + copies the snapshot version over target_path
    - returns error when the live path is not writable
    # restore
    -> std.fs.read_all
    -> std.fs.write_all
