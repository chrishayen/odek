# Requirement: "a searchable collection of code snippets"

A library that stores snippets keyed by name, tagged for discovery, and supports full-text search across titles and bodies.

std: (all units exist)

snippet_box
  snippet_box.new
    @ () -> box_state
    + creates an empty snippet collection
    # construction
  snippet_box.add
    @ (state: box_state, name: string, body: string, tags: list[string]) -> result[box_state, string]
    + adds a snippet with the given tags
    - returns error when a snippet with the same name already exists
    # registration
  snippet_box.get
    @ (state: box_state, name: string) -> result[snippet, string]
    + returns a snippet by name
    - returns error when the snippet is unknown
    # lookup
  snippet_box.find_by_tag
    @ (state: box_state, tag: string) -> list[snippet]
    + returns all snippets with the given tag
    # lookup
  snippet_box.search
    @ (state: box_state, query: string) -> list[snippet]
    + returns snippets whose name or body contains the query (case-insensitive)
    # search
