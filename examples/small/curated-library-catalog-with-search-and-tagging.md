# Requirement: "a curated library catalog with search and tagging"

A small catalog data structure with entries that can be registered, queried by tag, and searched by name.

std: (all units exist)

catalog
  catalog.new
    @ () -> catalog_state
    + returns an empty catalog
    # construction
  catalog.add_entry
    @ (cat: catalog_state, name: string, description: string, tags: list[string]) -> catalog_state
    + adds a new entry and indexes it under each of its tags
    ? duplicate names replace the prior entry
    # registration
  catalog.find_by_tag
    @ (cat: catalog_state, tag: string) -> list[entry]
    + returns entries tagged with tag
    + returns an empty list when the tag is unknown
    # query
  catalog.search_name
    @ (cat: catalog_state, query: string) -> list[entry]
    + returns entries whose name or description contains query, case-insensitive
    # search
