# Requirement: "a runtime version management library"

Tracks installed runtime versions, resolves the active version using per-directory and global pins.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads the entire file as text
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes text, replacing any existing file
      # filesystem
    std.fs.exists
      @ (path: string) -> bool
      + returns true when a file or directory exists at the path
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns the names of entries in the directory
      - returns error when the directory cannot be read
      # filesystem
  std.path
    std.path.join
      @ (parts: list[string]) -> string
      + joins path segments with the platform separator
      # paths
    std.path.parent
      @ (path: string) -> string
      + returns the parent directory of the given path
      # paths

version_mgr
  version_mgr.list_installed
    @ (root_dir: string) -> result[list[string], string]
    + returns the names of subdirectories under the versions root
    - returns error when the root does not exist
    # inventory
    -> std.fs.list_dir
    -> std.path.join
  version_mgr.is_installed
    @ (root_dir: string, version: string) -> bool
    + returns true when a directory for the version exists
    # inventory
    -> std.path.join
    -> std.fs.exists
  version_mgr.set_global
    @ (config_dir: string, version: string) -> result[void, string]
    + writes the version to the global pin file
    # pinning
    -> std.path.join
    -> std.fs.write_all
  version_mgr.set_local
    @ (dir: string, version: string) -> result[void, string]
    + writes a version pin in the given directory
    # pinning
    -> std.path.join
    -> std.fs.write_all
  version_mgr.resolve
    @ (start_dir: string, config_dir: string) -> result[string, string]
    + walks up from start_dir looking for a local pin and falls back to the global pin
    - returns error when no pin is found anywhere
    # resolution
    -> std.path.join
    -> std.path.parent
    -> std.fs.exists
    -> std.fs.read_all
