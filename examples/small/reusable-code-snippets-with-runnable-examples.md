# Requirement: "a library of reusable code snippets with runnable examples"

A library that stores named snippets with descriptions and expected outputs, and runs a snippet's example to compare actual output against the expected string.

std: (all units exist)

snippets
  snippets.new_catalog
    @ () -> catalog_state
    + creates an empty snippet catalog
    # construction
  snippets.add
    @ (state: catalog_state, name: string, description: string, expected_output: string) -> result[catalog_state, string]
    + registers a snippet with a description and expected output
    - returns error when a snippet with the same name already exists
    # registration
  snippets.get
    @ (state: catalog_state, name: string) -> result[snippet_entry, string]
    + returns the snippet entry for a name
    - returns error when the snippet is unknown
    # lookup
  snippets.check_example
    @ (entry: snippet_entry, actual_output: string) -> bool
    + returns true when the actual output matches the snippet's expected output
    # verification
