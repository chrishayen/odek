# Requirement: "a programming-language reference lookup"

An indexed reference of language topics with fuzzy search and cross-links between entries.

std: (all units exist)

langref
  langref.new_index
    @ () -> index_state
    + creates an empty reference index
    # construction
  langref.add_topic
    @ (state: index_state, id: string, title: string, body: string, tags: list[string]) -> index_state
    + adds or replaces a topic keyed by id
    # ingest
  langref.add_cross_link
    @ (state: index_state, from_id: string, to_id: string) -> result[index_state, string]
    + records a directed link between two topics
    - returns error when either topic id does not exist
    # graph
  langref.get_topic
    @ (state: index_state, id: string) -> optional[topic]
    + returns the topic if present
    - returns none when the id is unknown
    # read
  langref.search
    @ (state: index_state, query: string, limit: i32) -> list[search_hit]
    + returns topics ranked by substring matches in title (weighted 3x) and body (weighted 1x)
    + returns at most limit hits
    - returns an empty list when no topic matches
    # search
  langref.related
    @ (state: index_state, id: string) -> list[string]
    + returns the ids linked from the given topic in insertion order
    # graph
