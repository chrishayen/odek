# Requirement: "a web framework that generates an OpenAPI 3 specification from route definitions"

Routes are declared with typed request and response schemas. The framework stores them, dispatches requests, and emits an OpenAPI 3 document describing the whole API.

std
  std.json
    std.json.encode_object
      fn (obj: map[string,string]) -> string
      + encodes a string map as a JSON object
      # serialization
  std.collections
    std.collections.map_keys_sorted
      fn (m: map[string,string]) -> list[string]
      + returns keys in lexicographic order for deterministic output
      # collections

webapi
  webapi.new
    fn (title: string, version: string) -> api_state
    + creates an empty API with the given title and version
    # construction
  webapi.add_route
    fn (state: api_state, method: string, path: string, request_schema: string, response_schema: string, handler_id: string) -> api_state
    + stores a typed route with references to named schemas
    - returns unchanged state when (method, path) is already registered
    # route_registration
  webapi.add_schema
    fn (state: api_state, name: string, json_schema: string) -> api_state
    + stores a reusable schema under the given name
    # schema_registry
  webapi.dispatch
    fn (state: api_state, method: string, path: string) -> result[string, string]
    + returns the handler_id for (method, path)
    - returns error "not found" when no matching route exists
    # routing
  webapi.to_openapi_json
    fn (state: api_state) -> string
    + emits an OpenAPI 3 document with info, paths, and components.schemas
    ? paths and schema names are walked in sorted order so output is deterministic
    # spec_generation
    -> std.collections.map_keys_sorted
    -> std.json.encode_object
