# Requirement: "a library that turns a command-line argument spec into a GUI form"

Given an argument spec, produces a form definition and serializes user input back into a command line. Rendering the form is the caller's job.

std: (all units exist)

cli_to_gui
  cli_to_gui.parse_spec
    @ (raw: string) -> result[list[arg_field], string]
    + parses a simple "--name:type:help" line-based spec into form fields
    - returns error on an unrecognized type token
    # spec_parsing
  cli_to_gui.build_form
    @ (fields: list[arg_field]) -> form
    + returns a form with one widget per field: text, number, bool, or choice
    # form_construction
  cli_to_gui.to_command_line
    @ (fields: list[arg_field], values: map[string, string]) -> list[string]
    + renders values as argv tokens honoring each field's flag style
    + omits boolean flags whose value is false
    # serialization
  cli_to_gui.validate
    @ (fields: list[arg_field], values: map[string, string]) -> list[string]
    + returns one error message per required field that is blank or malformed
    # validation
