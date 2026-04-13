# Requirement: "a library for looking up entries in themed trivia datasets"

Pure in-memory lookup into named datasets the caller registers. The library does not ship any dataset content.

std: (all units exist)

trivia
  trivia.new
    @ () -> trivia_state
    + creates an empty store containing zero datasets
    # construction
  trivia.register_dataset
    @ (state: trivia_state, topic: string, entries: map[string, string]) -> trivia_state
    + stores the dataset under the given topic key
    ? existing datasets under the same topic are replaced
    # registration
  trivia.get
    @ (state: trivia_state, topic: string, key: string) -> result[string, string]
    + returns the entry for (topic, key)
    - returns error when topic is unknown
    - returns error when key is unknown within the topic
    # lookup
  trivia.list_keys
    @ (state: trivia_state, topic: string) -> result[list[string], string]
    + returns all keys in the dataset for the topic, sorted lexicographically
    - returns error when topic is unknown
    # introspection
  trivia.topics
    @ (state: trivia_state) -> list[string]
    + returns all registered topic names, sorted lexicographically
    # introspection
