# Requirement: "an application scaffold applying modern practices"

A library that assembles a runnable application from configuration, logging, and a health probe.

std
  std.env
    std.env.read_var
      fn (name: string) -> optional[string]
      + returns the value of an environment variable when set
      - returns none when the variable is missing
      # environment

app_scaffold
  app_scaffold.load_config
    fn (prefix: string) -> result[map[string, string], string]
    + collects environment variables matching the prefix into a config map
    - returns error when a required key is missing
    # configuration
    -> std.env.read_var
  app_scaffold.build
    fn (config: map[string, string]) -> result[app_state, string]
    + wires a logger and a health probe into an application state
    - returns error when config is missing the log level
    # wiring
  app_scaffold.health
    fn (state: app_state) -> bool
    + returns true when the application is ready to serve
    - returns false before build has completed
    # health
