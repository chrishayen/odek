# Requirement: "a library for managing multiple installed versions of a runtime"

Tracks installed versions, downloads new ones, and switches the active version by updating a symlink.

std
  std.http
    std.http.download
      fn (url: string, dest_path: string) -> result[void, string]
      + downloads a URL to the given filesystem path
      - returns error on network failure
      - returns error on filesystem failure
      # http
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns the entries of a directory
      - returns error when the path is not a directory
      # filesystem
    std.fs.symlink
      fn (target: string, link_path: string) -> result[void, string]
      + creates or replaces a symlink at link_path pointing at target
      - returns error on filesystem failure
      # filesystem
    std.fs.remove_dir
      fn (path: string) -> result[void, string]
      + recursively removes a directory
      - returns error on filesystem failure
      # filesystem

version_manager
  version_manager.list_installed
    fn (install_root: string) -> result[list[string], string]
    + returns the version names present under install_root
    - returns error when install_root is missing
    # inventory
    -> std.fs.list_dir
  version_manager.install
    fn (install_root: string, version: string, download_url: string) -> result[void, string]
    + downloads the given version's archive into install_root/version/
    - returns error on network failure
    - returns error when the version is already installed
    # install
    -> std.http.download
  version_manager.use
    fn (install_root: string, version: string, active_link: string) -> result[void, string]
    + points active_link at install_root/version/
    - returns error when the requested version is not installed
    # switching
    -> std.fs.symlink
  version_manager.uninstall
    fn (install_root: string, version: string) -> result[void, string]
    + removes install_root/version/ from disk
    - returns error when the version is not installed
    # uninstall
    -> std.fs.remove_dir
  version_manager.current
    fn (active_link: string) -> result[string, string]
    + returns the version name the active link currently points at
    - returns error when no active link exists
    # inspection
