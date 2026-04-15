# Requirement: "a library that builds and runs tests across multiple isolated environments"

Reads a config describing environments, provisions an isolated directory per environment, installs dependencies, and runs commands.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads file contents as text
      - returns error when missing
      # io
    std.fs.make_dir_all
      fn (path: string) -> result[void, string]
      + creates a directory and any missing parents
      # io
    std.fs.remove_dir_all
      fn (path: string) -> result[void, string]
      + recursively removes a directory
      # io
  std.process
    std.process.run
      fn (program: string, args: list[string], cwd: string, env: map[string,string]) -> result[process_result, string]
      + runs a subprocess and returns stdout, stderr, and exit code
      - returns error when the program cannot be launched
      # process
  std.text
    std.text.parse_ini
      fn (text: string) -> result[map[string, map[string,string]], string]
      + parses an INI-style config into sections and keys
      - returns error on malformed headers
      # parsing

tox
  tox.config_load
    fn (path: string) -> result[tox_config, string]
    + reads and parses the environment definitions
    # config
    -> std.fs.read_all
    -> std.text.parse_ini
  tox.list_envs
    fn (config: tox_config) -> list[string]
    + returns the declared environment names
    # config
  tox.env_create
    fn (config: tox_config, name: string, workdir: string) -> result[env_handle, string]
    + provisions an isolated directory for the environment
    - returns error when the environment is not declared
    # provisioning
    -> std.fs.make_dir_all
  tox.env_install
    fn (env: env_handle, deps: list[string]) -> result[void, string]
    + runs the dependency installer inside the environment
    # provisioning
    -> std.process.run
  tox.env_run_command
    fn (env: env_handle, program: string, args: list[string]) -> result[process_result, string]
    + executes a command within the environment
    - returns error on nonzero exit
    # execution
    -> std.process.run
  tox.env_destroy
    fn (env: env_handle) -> result[void, string]
    + removes the provisioned directory
    # cleanup
    -> std.fs.remove_dir_all
  tox.run_all
    fn (config: tox_config, workdir: string) -> result[list[env_result], string]
    + provisions, installs, and runs commands for every declared environment
    - returns error when any environment fails and strict mode is on
    # orchestration
