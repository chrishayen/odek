# Requirement: "a functional-options configuration builder"

A tiny builder pattern: start with defaults, apply a list of option functions, get the final config.

std: (all units exist)

options
  options.default_config
    fn () -> config_state
    + returns a config populated with default values for every field
    # defaults
  options.apply
    fn (base: config_state, opts: list[func(config_state) -> config_state]) -> config_state
    + folds the option list over the base config in order
    + returns base unchanged when opts is empty
    # composition
  options.with_field
    fn (name: string, value: string) -> func(config_state) -> config_state
    + returns an option that sets the named field to the given value
    ? unknown field names are ignored by the returned option
    # option
