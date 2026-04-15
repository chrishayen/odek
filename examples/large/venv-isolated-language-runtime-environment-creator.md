# Requirement: "a library for creating isolated language runtime environments"

Creates a directory that holds a copy (or symlink) of an interpreter, a local package store, and an activation script that prepends the environment's binary directory to PATH.

std
  std.fs
    std.fs.make_dir
      fn (path: string) -> result[void, string]
      + creates the directory, including missing parents
      - returns error when path exists as a file
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to the file, creating or overwriting
      - returns error when the parent directory is missing
      # filesystem
    std.fs.symlink
      fn (target: string, link_path: string) -> result[void, string]
      + creates a symbolic link at link_path pointing to target
      - returns error when link_path already exists
      # filesystem
    std.fs.path_exists
      fn (path: string) -> bool
      + returns true when the path exists
      # filesystem
  std.os
    std.os.run_command
      fn (program: string, args: list[string], env: map[string, string]) -> result[process_result, string]
      + runs the program and returns exit code, stdout, and stderr
      - returns error when the program cannot be spawned
      # process

venv
  venv.layout
    fn (root: string) -> venv_paths
    + returns the directories a new environment should have (bin, lib, include, site-packages)
    # layout
  venv.create
    fn (root: string, interpreter: string) -> result[venv_paths, string]
    + creates the directory layout and links the interpreter into the bin directory
    - returns error when root already exists
    - returns error when the interpreter is not found on disk
    # creation
    -> venv.layout
    -> std.fs.path_exists
    -> std.fs.make_dir
    -> std.fs.symlink
  venv.write_activate_script
    fn (paths: venv_paths) -> result[void, string]
    + writes a POSIX shell activation script that prepends bin to PATH and sets a prompt prefix
    # activation
    -> std.fs.write_all
  venv.write_config
    fn (paths: venv_paths, interpreter: string, version: string) -> result[void, string]
    + writes a config file recording the interpreter path and version
    # metadata
    -> std.fs.write_all
  venv.install_package
    fn (paths: venv_paths, package_name: string) -> result[void, string]
    + runs the package manager inside the environment to install the named package
    - returns error when the install process exits non-zero
    # package_management
    -> std.os.run_command
  venv.bootstrap
    fn (root: string, interpreter: string) -> result[venv_paths, string]
    + creates the layout, writes the activation script, and writes the config
    - returns error from whichever step fails first
    # orchestration
    -> venv.create
    -> venv.write_activate_script
    -> venv.write_config
