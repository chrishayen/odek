# Requirement: "an all-in-one toolkit for a scripting language runtime"

A toolkit combines a runtime entry, a dependency resolver, a script executor, a test runner, and a bundler. Each subsystem gets a thin project face backed by substantive std primitives.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns file contents as a string
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: string) -> result[void, string]
      + writes data to path, overwriting any existing file
      - returns error when the parent directory does not exist
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns all file paths rooted at root in depth-first order
      - returns error when root is not a directory
      # filesystem
  std.process
    std.process.spawn
      @ (argv: list[string], cwd: string) -> result[i32, string]
      + runs a child process and returns its exit code
      - returns error when the executable is not found
      # process
  std.json
    std.json.parse
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into string keyed entries
      - returns error on invalid JSON
      # serialization
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex digest of the input
      # hashing

toolkit
  toolkit.run_script
    @ (script_path: string, argv: list[string]) -> result[i32, string]
    + loads a script file and executes it, returning the exit code
    - returns error when the script file cannot be read
    # script_execution
    -> std.fs.read_all
    -> std.process.spawn
  toolkit.resolve_dependencies
    @ (manifest_path: string) -> result[map[string, string], string]
    + reads the manifest and returns a name-to-version map of direct dependencies
    - returns error on malformed manifest
    # dependency_resolution
    -> std.fs.read_all
    -> std.json.parse
  toolkit.install_dependency
    @ (name: string, version: string, target_dir: string) -> result[string, string]
    + downloads the named dependency and writes it under target_dir, returning the install path
    - returns error when the target directory is not writable
    # package_install
    -> std.fs.write_all
    -> std.hash.sha256_hex
  toolkit.run_tests
    @ (root: string) -> result[test_report, string]
    + discovers test files under root and returns a report with pass/fail counts
    - returns error when root contains no test files
    # test_runner
    -> std.fs.walk
    -> std.process.spawn
  toolkit.bundle
    @ (entry: string, output: string) -> result[void, string]
    + resolves the dependency graph from entry and writes a single-file bundle to output
    - returns error when an import cannot be resolved
    # bundler
    -> std.fs.read_all
    -> std.fs.write_all
  toolkit.format_source
    @ (source: string) -> string
    + returns a canonically formatted version of the input source
    # formatting
  toolkit.type_check
    @ (source: string) -> result[void, list[string]]
    + returns ok when the source type-checks
    - returns a list of error messages when it does not
    # type_checking
