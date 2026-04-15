# Requirement: "a library that trims, sanitizes, and scrubs user input based on field tags"

Declaratively apply transforms to named fields. Each transform is a small, standalone rune.

std: (all units exist)

conform
  conform.register_rule
    fn (state: conform_state, field: string, transform: string) -> conform_state
    + registers a named transform to run on a field
    ? transforms are applied in registration order
    # registration
  conform.apply
    fn (state: conform_state, input: map[string, string]) -> map[string, string]
    + returns a new map with every field transformed per its registered rules
    # application
  conform.trim
    fn (value: string) -> string
    + removes leading and trailing whitespace
    # transform
  conform.lowercase
    fn (value: string) -> string
    + returns the value with all ASCII letters lowered
    # transform
  conform.uppercase
    fn (value: string) -> string
    + returns the value with all ASCII letters raised
    # transform
  conform.collapse_spaces
    fn (value: string) -> string
    + replaces runs of whitespace with a single space
    # transform
  conform.strip_html
    fn (value: string) -> string
    + removes anything that looks like an HTML tag
    # transform
  conform.strip_non_numeric
    fn (value: string) -> string
    + removes every character that is not a digit
    # transform
