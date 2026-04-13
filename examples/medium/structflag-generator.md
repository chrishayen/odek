# Requirement: "a struct-based command-line flag generator"

Describe a set of options as a schema; generate a flag spec and parse argv against it to fill a typed value map.

std: (all units exist)

structflags
  structflags.new_schema
    @ () -> schema_state
    + creates an empty flag schema
    # construction
  structflags.add_string
    @ (schema: schema_state, name: string, default_value: string, help: string) -> schema_state
    + adds a string-valued flag
    # schema
  structflags.add_int
    @ (schema: schema_state, name: string, default_value: i64, help: string) -> schema_state
    + adds an integer-valued flag
    # schema
  structflags.add_bool
    @ (schema: schema_state, name: string, default_value: bool, help: string) -> schema_state
    + adds a boolean flag (presence implies true unless value provided)
    # schema
  structflags.parse_argv
    @ (schema: schema_state, argv: list[string]) -> result[parsed_flags, string]
    + accepts --name=value, --name value, and --flag (bool) forms
    + applies defaults for flags that were not supplied
    - returns error on unknown flag names
    - returns error when a required value is missing or wrong type
    # parsing
  structflags.get_string
    @ (parsed: parsed_flags, name: string) -> result[string, string]
    + returns the parsed string value for a flag
    - returns error when the flag was not declared as a string
    # access
  structflags.get_int
    @ (parsed: parsed_flags, name: string) -> result[i64, string]
    + returns the parsed integer value for a flag
    - returns error when the flag was not declared as an integer
    # access
  structflags.get_bool
    @ (parsed: parsed_flags, name: string) -> result[bool, string]
    + returns the parsed boolean value for a flag
    - returns error when the flag was not declared as a boolean
    # access
  structflags.format_help
    @ (schema: schema_state) -> string
    + renders a multi-line help listing with names, types, defaults, and descriptions
    # help
