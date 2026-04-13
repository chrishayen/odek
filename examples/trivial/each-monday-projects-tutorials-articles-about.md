# Requirement: "a curated weekly newsletter item store"

The source description is a newsletter, not a library. Interpret as a tiny store that returns the current week's curated items.

std: (all units exist)

weekly_digest
  weekly_digest.items_for_week
    @ (week_number: i32) -> list[string]
    + returns the list of curated items for the given ISO week number
    - returns an empty list when no items are recorded for that week
    ? items are kept in an in-memory map keyed by week number
    # retrieval
