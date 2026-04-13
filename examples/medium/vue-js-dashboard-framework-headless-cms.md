# Requirement: "a headless content management library for dashboards"

Stores content entries organized by type and exposes a simple query interface. No UI, no HTTP layer — pure in-memory store with pluggable persistence.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.generate
      @ () -> string
      + returns a new unique identifier
      # id_generation

cms
  cms.new_store
    @ () -> cms_state
    + returns an empty content store
    # construction
  cms.define_type
    @ (state: cms_state, type_name: string, fields: list[string]) -> result[cms_state, string]
    + registers a new content type with the given field names
    - returns error when the type already exists
    # schema
  cms.create_entry
    @ (state: cms_state, type_name: string, data: map[string, string]) -> result[tuple[string, cms_state], string]
    + creates an entry and returns its new id with the updated state
    - returns error when the type is unknown
    - returns error when required fields are missing
    # write
    -> std.id.generate
    -> std.time.now_millis
  cms.get_entry
    @ (state: cms_state, id: string) -> optional[cms_entry]
    + returns the entry when present
    - returns none when the id is unknown
    # read
  cms.update_entry
    @ (state: cms_state, id: string, data: map[string, string]) -> result[cms_state, string]
    + merges the new data into the existing entry
    - returns error when the id is unknown
    # write
    -> std.time.now_millis
  cms.delete_entry
    @ (state: cms_state, id: string) -> result[cms_state, string]
    + removes the entry
    - returns error when the id is unknown
    # write
  cms.list_by_type
    @ (state: cms_state, type_name: string) -> list[cms_entry]
    + returns all entries of the given type in creation order
    # query
