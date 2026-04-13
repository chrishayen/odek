# Requirement: "a library that exposes a curated collection of tutorial entries"

Tutorials are static data with a simple lookup by topic.

std: (all units exist)

tutorials
  tutorials.list_topics
    @ () -> list[string]
    + returns all available tutorial topic names
    # catalog
  tutorials.get
    @ (topic: string) -> optional[string]
    + returns the tutorial body for a known topic
    - returns none for unknown topics
    # lookup
