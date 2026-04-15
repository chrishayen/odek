# Requirement: "a portable development environment manager"

Installs, registers, and activates self-contained toolchain bundles under a portable root directory.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads the full file as UTF-8
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes contents to path, overwriting
      - returns error when the parent directory is missing
      # filesystem

portdev
  portdev.init_root
    fn (root: string) -> result[env_state, string]
    + creates the layout under root and returns an environment handle
    - returns error when root exists and is not empty
    # construction
  portdev.install_bundle
    fn (state: env_state, bundle_path: string) -> result[env_state, string]
    + extracts a toolchain bundle into the environment and registers it
    - returns error when the bundle manifest is missing
    # installation
    -> std.fs.read_all
  portdev.list_bundles
    fn (state: env_state) -> list[bundle_info]
    + returns the installed bundles with name and version
    # introspection
  portdev.build_activation_script
    fn (state: env_state) -> string
    + returns a shell script that prepends bundle paths to PATH and sets tool variables
    # activation
    -> std.fs.write_all
