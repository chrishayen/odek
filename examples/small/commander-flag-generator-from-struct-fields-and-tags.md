# Requirement: "a library that derives command-line flag definitions from struct fields and tags"

Introspects a struct-like type description and builds a flag parser with defaults and usage text.

std: (all units exist)

commandeer
  commandeer.describe_struct
    fn (type_name: string, fields: list[field_meta]) -> struct_desc
    + records field names, types, default values, and tag metadata
    # reflection
  commandeer.build_parser
    fn (desc: struct_desc) -> flag_parser
    + creates a flag parser whose options mirror the struct's fields
    # construction
  commandeer.parse_args
    fn (parser: flag_parser, args: list[string]) -> result[map[string, string], string]
    + parses the argument list into field values, applying defaults for missing flags
    - returns error on unknown flag
    - returns error on missing value for a required flag
    # parsing
  commandeer.render_usage
    fn (parser: flag_parser) -> string
    + produces a human-readable usage block listing every flag with its default
    # usage
