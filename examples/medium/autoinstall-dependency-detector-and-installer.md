# Requirement: "a library that detects missing dependencies in source files and installs them"

Scans source files for import-like statements, diffs against the declared manifest, and runs an install command for anything missing.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file as text
      - returns error when missing
      # io
    std.fs.write_all
      fn (path: string, content: string) -> result[void, string]
      + writes text to a file
      # io
    std.fs.watch
      fn (path: string) -> result[watch_state, string]
      + subscribes to change events under a directory
      - returns error when the path is not watchable
      # io
    std.fs.next_event
      fn (watch: watch_state) -> result[fs_event, string]
      + blocks for the next change event
      # io
  std.process
    std.process.run
      fn (program: string, args: list[string], cwd: string, env: map[string,string]) -> result[process_result, string]
      + runs a subprocess and returns its output
      - returns error when the program cannot be launched
      # process

autoinstall
  autoinstall.scan_imports
    fn (source: string) -> list[string]
    + extracts import specifiers from a source file
    # scanning
  autoinstall.manifest_load
    fn (path: string) -> result[map[string,string], string]
    + reads the project's declared dependencies
    - returns error when the manifest is missing
    # manifest
    -> std.fs.read_all
  autoinstall.manifest_save
    fn (path: string, deps: map[string,string]) -> result[void, string]
    + writes the dependency map back to the manifest
    # manifest
    -> std.fs.write_all
  autoinstall.missing
    fn (imports: list[string], manifest: map[string,string]) -> list[string]
    + returns import specifiers not present in the manifest
    # diff
  autoinstall.install
    fn (dep: string, workdir: string) -> result[void, string]
    + invokes the install command for a single dependency
    - returns error on nonzero exit
    # install
    -> std.process.run
  autoinstall.sync_file
    fn (source_path: string, manifest_path: string, workdir: string) -> result[list[string], string]
    + scans, diffs, installs, and returns the dependencies that were added
    # orchestration
    -> std.fs.read_all
  autoinstall.watch_loop
    fn (project_dir: string, manifest_path: string) -> result[void, string]
    + watches for source changes and re-syncs dependencies on each edit
    # watcher
    -> std.fs.watch
    -> std.fs.next_event
