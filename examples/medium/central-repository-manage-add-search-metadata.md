# Requirement: "a central repository to manage (add, search, and query metadata of) one-liners, scripts, and tools"

A persistent catalog of script snippets with tags and free-text search. Storage is a thin std primitive so the project layer focuses on the catalog semantics.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads entire file contents
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes atomically to path, creating parent directories
      - returns error on permission failure
      # filesystem
  std.json
    std.json.parse
      @ (raw: string) -> result[map[string, string], string]
      + parses a flat JSON object
      - returns error on malformed input
      # serialization
    std.json.encode
      @ (obj: map[string, string]) -> string
      + encodes a flat map as JSON
      # serialization

script_catalog
  script_catalog.load
    @ (path: string) -> result[catalog_state, string]
    + loads the catalog from a file, returning an empty catalog when the file is absent
    - returns error when the file exists but is malformed
    # persistence
    -> std.fs.read_all
    -> std.json.parse
  script_catalog.save
    @ (state: catalog_state, path: string) -> result[void, string]
    + writes the catalog to disk
    # persistence
    -> std.fs.write_all
    -> std.json.encode
  script_catalog.add
    @ (state: catalog_state, alias: string, body: string, tags: list[string]) -> result[catalog_state, string]
    + registers a new entry under a unique alias
    - returns error when alias already exists
    - returns error when alias is empty
    # mutation
  script_catalog.remove
    @ (state: catalog_state, alias: string) -> result[catalog_state, string]
    + removes the entry for the alias
    - returns error when alias is not present
    # mutation
  script_catalog.get
    @ (state: catalog_state, alias: string) -> optional[script_entry]
    + returns the entry when present
    - returns none when alias is unknown
    # lookup
  script_catalog.search
    @ (state: catalog_state, query: string) -> list[script_entry]
    + returns entries whose alias, body, or tags contain the query substring
    + returns empty list on no match
    ? search is case-insensitive substring match
    # search
  script_catalog.list_by_tag
    @ (state: catalog_state, tag: string) -> list[script_entry]
    + returns every entry that carries the tag
    # search
