# Requirement: "a library for updating installed binaries by re-running their original install command"

Scans an install directory, extracts module metadata embedded in each binary, checks the registry for a newer version, and reinstalls outdated ones.

std
  std.fs
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns file names in the directory
      - returns error when the path does not exist
      # filesystem
    std.fs.is_executable
      @ (path: string) -> bool
      + true when the file has the executable bit set
      # filesystem
  std.process
    std.process.run
      @ (cmd: string, args: list[string]) -> result[string, string]
      + runs a subprocess and returns its captured stdout
      - returns error with stderr when the exit code is non-zero
      # process
  std.http
    std.http.get_json
      @ (url: string) -> result[map[string, string], string]
      + performs a GET and parses the JSON body into a string-to-string map
      - returns error on non-2xx status or invalid JSON
      # network

bin_updater
  bin_updater.scan_install_dir
    @ (dir: string) -> result[list[string], string]
    + returns absolute paths of executable files in the directory
    - returns error when the directory is unreadable
    # discovery
    -> std.fs.list_dir
    -> std.fs.is_executable
  bin_updater.read_binary_metadata
    @ (path: string) -> result[binary_metadata, string]
    + extracts module path and version string embedded in the binary
    - returns error when no metadata is present
    # inspection
    -> std.process.run
  bin_updater.fetch_latest_version
    @ (module_path: string) -> result[string, string]
    + returns the latest published version for the given module
    - returns error when the module is unknown to the registry
    # registry
    -> std.http.get_json
  bin_updater.is_outdated
    @ (current: string, latest: string) -> bool
    + true when latest compares greater than current under semver ordering
    - false when versions are equal
    # comparison
  bin_updater.reinstall
    @ (module_path: string, version: string) -> result[void, string]
    + invokes the install tool to install the requested module at the given version
    - returns error when the install command fails
    # installation
    -> std.process.run
  bin_updater.update_all
    @ (dir: string) -> result[list[string], string]
    + returns the module paths that were reinstalled
    # orchestration
