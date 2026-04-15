# Requirement: "a CLI framework that builds commands from typed function signatures"

Users declare commands with typed parameters, then dispatch an argv array to the right command.

std: (all units exist)

cli
  cli.new_app
    fn (name: string) -> app_state
    + creates an empty application with the given program name
    # construction
  cli.command
    fn (app: app_state, name: string, params: list[param_spec], handler: fn(map[string, value]) -> i32) -> app_state
    + registers a named command with typed parameters and a handler returning an exit code
    ? each param_spec has a name, a type (string, int, float, bool), a kind (positional, flag, option), and an optional default
    # registration
  cli.parse
    fn (app: app_state, argv: list[string]) -> result[tuple[string, map[string, value]], string]
    + resolves the command name and a map of typed parameter values
    - returns error on unknown command, unknown flag, missing required positional, or type mismatch
    # parsing
  cli.dispatch
    fn (app: app_state, argv: list[string]) -> i32
    + parses argv and invokes the matched handler, returning its exit code
    - returns 2 with an error-formatted usage on parse failure
    # dispatch
  cli.help
    fn (app: app_state, command: optional[string]) -> string
    + renders usage text for the app or a specific command
    # help
