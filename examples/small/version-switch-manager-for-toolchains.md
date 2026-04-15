# Requirement: "a version manager that lets a user switch between installed toolchain versions"

Tracks installed versions on disk and maintains a pointer to the currently active one.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns the immediate entries of a directory
      - returns error when the directory does not exist
      # filesystem
    std.fs.symlink
      fn (target: string, link_path: string) -> result[void, string]
      + creates or replaces a symbolic link pointing at target
      # filesystem
    std.fs.read_link
      fn (link_path: string) -> result[string, string]
      + returns the target path of a symbolic link
      - returns error when the path is not a link
      # filesystem

version_switch
  version_switch.list_installed
    fn (install_root: string) -> result[list[string], string]
    + returns every version directory under install_root
    # query
    -> std.fs.list_dir
  version_switch.current
    fn (current_link: string) -> result[string, string]
    + returns the version currently pointed to by the active link
    - returns error when no version is active
    # query
    -> std.fs.read_link
  version_switch.activate
    fn (install_root: string, version: string, current_link: string) -> result[void, string]
    + points current_link at install_root/version
    - returns error when the version is not installed
    # mutation
    -> std.fs.list_dir
    -> std.fs.symlink
