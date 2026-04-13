# Requirement: "a searchable store for shell command snippets"

Snippets live in a local store; users add, search by keyword, and retrieve by id.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error when path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data atomically, creating the file if needed
      - returns error when the parent directory does not exist
      # filesystem

snippet_store
  snippet_store.load
    @ (path: string) -> result[snippet_store_state, string]
    + reads snippets from a json-backed store file
    + returns an empty store when the file does not exist
    - returns error on malformed data
    # persistence
    -> std.fs.read_all
  snippet_store.save
    @ (state: snippet_store_state, path: string) -> result[void, string]
    + writes the store to path
    # persistence
    -> std.fs.write_all
  snippet_store.add
    @ (state: snippet_store_state, command: string, description: string, tags: list[string]) -> tuple[string, snippet_store_state]
    + adds a snippet and returns its new id
    # mutation
  snippet_store.search
    @ (state: snippet_store_state, query: string) -> list[snippet]
    + returns snippets whose command, description, or tags contain query
    + ranks by number of matching terms
    - returns empty list when nothing matches
    # search
  snippet_store.get
    @ (state: snippet_store_state, id: string) -> result[snippet, string]
    + returns the snippet for an id
    - returns error when id is unknown
    # retrieval
  snippet_store.remove
    @ (state: snippet_store_state, id: string) -> result[snippet_store_state, string]
    + removes the snippet with the given id
    - returns error when id is unknown
    # mutation
