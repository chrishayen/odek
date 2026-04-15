# Requirement: "a headless content management framework"

Authors define content types; the framework stores typed entries and exposes them through a query API. Storage is abstracted behind a store handle.

std
  std.uuid
    std.uuid.new_v4
      fn () -> string
      + returns a random UUID string
      # identifiers
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.json
    std.json.encode_value
      fn (value: json_value) -> string
      + encodes a generic JSON value to a string
      # serialization
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document
      - returns error on invalid input
      # serialization

cms
  cms.define_type
    fn (schema: schema_registry, name: string, fields: list[field_def]) -> result[schema_registry, string]
    + registers a content type with named, typed fields
    - returns error when a type with that name already exists
    # schema
  cms.validate_entry
    fn (schema: schema_registry, type_name: string, values: map[string, json_value]) -> result[void, list[string]]
    + returns ok when all required fields are present and each value matches its declared type
    - returns a list of per-field errors otherwise
    # validation
  cms.create_entry
    fn (store: store_handle, schema: schema_registry, type_name: string, values: map[string, json_value]) -> result[string, string]
    + validates the entry and persists it, returning the new entry id
    - returns error when validation fails
    # authoring
    -> std.uuid.new_v4
    -> std.time.now_millis
    -> std.json.encode_value
  cms.update_entry
    fn (store: store_handle, schema: schema_registry, id: string, patch: map[string, json_value]) -> result[void, string]
    + applies the patch after validating the merged result
    - returns error when the id does not exist
    # authoring
    -> std.time.now_millis
  cms.delete_entry
    fn (store: store_handle, id: string) -> result[void, string]
    + removes the entry
    - returns error when the id does not exist
    # authoring
  cms.get_entry
    fn (store: store_handle, id: string) -> result[entry, string]
    + returns the entry with all fields
    - returns error when the id does not exist
    # retrieval
    -> std.json.parse
  cms.list_entries
    fn (store: store_handle, type_name: string, limit: i32, offset: i32) -> result[list[entry], string]
    + returns a page of entries for a content type
    # retrieval
  cms.query_by_field
    fn (store: store_handle, type_name: string, field: string, value: json_value) -> result[list[entry], string]
    + returns entries where the given field equals the given value
    - returns error when the field is not part of the type
    # query
