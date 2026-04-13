# Requirement: "a RESTful endpoint generator for ORM models"

Given a model descriptor, produce handler functions for the standard CRUD operations.

std
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
    std.json.encode_list
      @ (items: list[map[string, string]]) -> string
      + encodes a list of objects as a JSON array
      # serialization

rest_gen
  rest_gen.define_model
    @ (name: string, fields: list[string], pk: string) -> model_descriptor
    + builds a descriptor capturing the resource name, field list, and primary key
    - returns an empty descriptor when name is empty
    # model_definition
  rest_gen.handle_list
    @ (model: model_descriptor, rows: list[map[string, string]]) -> string
    + returns a JSON array body for a GET collection request
    # list_endpoint
    -> std.json.encode_list
  rest_gen.handle_get
    @ (model: model_descriptor, row: optional[map[string, string]]) -> result[string, string]
    + returns the encoded row when present
    - returns "not found" error when row is absent
    # get_endpoint
    -> std.json.encode_object
  rest_gen.handle_create
    @ (model: model_descriptor, body: string) -> result[map[string, string], string]
    + parses the body and returns the new record mapping
    - returns error on invalid body or missing required field
    # create_endpoint
    -> std.json.parse_object
  rest_gen.handle_update
    @ (model: model_descriptor, body: string, existing: map[string, string]) -> result[map[string, string], string]
    + merges parsed fields into the existing row
    - returns error when parsing fails
    # update_endpoint
    -> std.json.parse_object
  rest_gen.handle_delete
    @ (model: model_descriptor, existed: bool) -> result[void, string]
    + returns ok when the row was removed
    - returns "not found" error when existed is false
    # delete_endpoint
