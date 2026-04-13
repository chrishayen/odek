# Requirement: "a command-line application framework with configuration loading and dependency injection"

A root command type with subcommand registration, argument parsing, a layered config loader, and a tiny dependency container the commands can pull from.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the whole file into memory
      - returns error when the path does not exist
      # io
  std.env
    std.env.get
      @ (key: string) -> optional[string]
      + returns the value of an environment variable, if set
      # environment

app
  app.new
    @ (name: string, version: string) -> app_state
    + constructs an empty application with the given name and version
    # construction
  app.register
    @ (state: app_state, name: string, handler: command_fn) -> app_state
    + registers a subcommand with its handler
    # registration
  app.parse_args
    @ (state: app_state, args: list[string]) -> result[parsed_invocation, string]
    + resolves the subcommand and parses flags into a typed bag
    - returns error when the subcommand is unknown
    - returns error when a required flag is missing
    # parsing
  app.load_config
    @ (paths: list[string]) -> result[map[string, string], string]
    + merges config from the given paths, later entries overriding earlier ones
    + falls back to environment variables for missing keys
    - returns error when a listed file is malformed
    # configuration
    -> std.fs.read_all
    -> std.env.get
  app.container_new
    @ () -> container_state
    + constructs an empty dependency container
    # injection
  app.container_register
    @ (state: container_state, key: string, value: dependency) -> container_state
    + stores a dependency under the given key
    # injection
  app.container_resolve
    @ (state: container_state, key: string) -> optional[dependency]
    + returns the stored dependency for key, if any
    # injection
  app.run
    @ (state: app_state, container: container_state, args: list[string]) -> result[i32, string]
    + parses args, resolves dependencies, and invokes the selected command
    + returns the command's exit code
    - returns error when parsing fails before a command is dispatched
    # execution
    -> app.parse_args
