# Requirement: "a library that returns annotated example programs for hands-on learning"

A read-only catalog of example snippets keyed by title.

std: (all units exist)

example_catalog
  example_catalog.titles
    @ () -> list[string]
    + returns every example title in curriculum order
    # query
  example_catalog.lookup
    @ (title: string) -> result[annotated_example, string]
    + returns the source and annotations for the example with the given title
    - returns error when no example has that title
    # query
