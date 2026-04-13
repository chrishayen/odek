# Requirement: "a library for building composable shell aliases"

An alias is a named command template with parameter substitution. The library stores, looks up, and expands aliases.

std
  std.strings
    std.strings.replace_all
      @ (s: string, old: string, new: string) -> string
      + returns s with every occurrence of old replaced by new
      # strings

alias
  alias.new_store
    @ () -> alias_store
    + constructs an empty alias store
    # construction
  alias.define
    @ (store: alias_store, name: string, template: string) -> alias_store
    + registers a new alias with a template containing {0}, {1}, ... placeholders
    + overwrites any existing alias with the same name
    # registration
  alias.expand
    @ (store: alias_store, name: string, args: list[string]) -> result[string, string]
    + substitutes positional args into the template and returns the result
    - returns error when name is not defined
    - returns error when the template references a position outside args
    # expansion
    -> std.strings.replace_all
  alias.list_names
    @ (store: alias_store) -> list[string]
    + returns the sorted names of every defined alias
    # introspection
