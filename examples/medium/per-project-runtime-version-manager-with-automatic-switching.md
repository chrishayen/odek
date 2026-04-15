# Requirement: "a per-project runtime version manager with automatic version switching"

Resolves the active version by walking up from a working directory looking for a pinned-version file, falling back to a global default.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads an entire file as UTF-8 text
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, content: string) -> result[void, string]
      + writes content to path, creating or overwriting
      - returns error on filesystem failure
      # filesystem
    std.fs.exists
      fn (path: string) -> bool
      + returns true when a file or directory exists at path
      # filesystem
    std.fs.parent_dir
      fn (path: string) -> optional[string]
      + returns the parent directory path, or none when path is the filesystem root
      # filesystem

version_resolver
  version_resolver.find_version_file
    fn (start_dir: string, filename: string) -> optional[string]
    + walks from start_dir upward until a file named filename is found and returns its absolute path
    + returns none when no ancestor directory contains the file
    # discovery
    -> std.fs.exists
    -> std.fs.parent_dir
  version_resolver.read_pinned_version
    fn (version_file_path: string) -> result[string, string]
    + returns the trimmed contents of the pinned-version file
    - returns error when the file cannot be read
    - returns error when the file is empty
    # reading
    -> std.fs.read_all
  version_resolver.pin_version
    fn (dir: string, filename: string, version: string) -> result[void, string]
    + writes the version into dir/filename
    - returns error on filesystem failure
    # writing
    -> std.fs.write_all
  version_resolver.resolve
    fn (start_dir: string, filename: string, global_default: string) -> string
    + returns the first version found by walking upward, falling back to global_default
    # resolution
    -> std.fs.exists
    -> std.fs.parent_dir
    -> std.fs.read_all
  version_resolver.set_global
    fn (global_config_path: string, version: string) -> result[void, string]
    + writes the global-default version to the global config file
    - returns error on filesystem failure
    # writing
    -> std.fs.write_all
