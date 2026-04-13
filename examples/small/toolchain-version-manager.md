# Requirement: "a library for managing multiple installed versions of a language toolchain"

std
  std.fs
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns the immediate entries of a directory
      - returns error when the directory does not exist
      # filesystem
    std.fs.write_text
      @ (path: string, content: string) -> result[void, string]
      + writes a text file, creating parents if needed
      # filesystem

toolchain_manager
  toolchain_manager.list_installed
    @ (root: string) -> result[list[string], string]
    + returns the version identifiers installed under the given root
    - returns error when the root directory is unreadable
    # inventory
    -> std.fs.list_dir
  toolchain_manager.install
    @ (root: string, version: string, archive: bytes) -> result[void, string]
    + unpacks the archive into the versioned directory
    - returns error when a matching version is already installed
    # install
  toolchain_manager.remove
    @ (root: string, version: string) -> result[void, string]
    + deletes the directory for the given version
    - returns error when the version is not installed
    # removal
  toolchain_manager.set_active
    @ (root: string, version: string) -> result[void, string]
    + records the given version as the active one
    - returns error when the version is not installed
    # selection
    -> std.fs.write_text
  toolchain_manager.get_active
    @ (root: string) -> result[optional[string], string]
    + returns the currently active version or none if unset
    # selection
