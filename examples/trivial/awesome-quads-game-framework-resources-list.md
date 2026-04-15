# Requirement: "a curated list of links to 2d game framework resources"

A curated list is static data. The library exposes the entries as a typed list.

std: (all units exist)

awesome_quads
  awesome_quads.list_entries
    fn () -> list[link_entry]
    + returns the full curated list of resource entries
    + each entry has title, url, and category fields
    ? data is compiled in; callers filter client-side
    # catalog
