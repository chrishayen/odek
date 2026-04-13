# Requirement: "a curated list of gpu api code and resources"

A static catalog of resources. One function returns the list.

std: (all units exist)

awesome_gpu
  awesome_gpu.list_entries
    @ () -> list[link_entry]
    + returns the full curated list of resource entries
    + each entry has title, url, and category fields
    ? data is compiled in; callers filter client-side
    # catalog
