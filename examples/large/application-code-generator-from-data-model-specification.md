# Requirement: "a library that generates application code from a data model specification"

Takes a model (entities, fields, relations) and emits source files for a CRUD backend and matching client stubs.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file as text
      - returns error when missing
      # io
    std.fs.write_all
      fn (path: string, content: string) -> result[void, string]
      + writes text to a file, creating parents as needed
      # io
    std.fs.make_dir_all
      fn (path: string) -> result[void, string]
      + creates a directory and missing parents
      # io
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses JSON into a value tree
      - returns error on malformed input
      # serialization
  std.template
    std.template.render
      fn (template: string, vars: map[string,string]) -> result[string, string]
      + substitutes {{name}} placeholders from the map
      - returns error on unclosed placeholders
      # templating

codegen
  codegen.spec_load
    fn (path: string) -> result[app_spec, string]
    + reads and parses a model specification file
    - returns error on invalid schema
    # loading
    -> std.fs.read_all
    -> std.json.parse
  codegen.spec_validate
    fn (spec: app_spec) -> result[void, string]
    + checks entities for name uniqueness and field type validity
    - returns error listing the first violation
    # validation
  codegen.entity_names
    fn (spec: app_spec) -> list[string]
    + returns the declared entity names
    # introspection
  codegen.emit_entity_model
    fn (spec: app_spec, entity: string) -> result[string, string]
    + renders the data model source for an entity
    - returns error when the entity is unknown
    # emission
    -> std.template.render
  codegen.emit_crud_routes
    fn (spec: app_spec, entity: string) -> result[string, string]
    + renders HTTP handlers for create, read, update, delete
    # emission
    -> std.template.render
  codegen.emit_client_stub
    fn (spec: app_spec, entity: string) -> result[string, string]
    + renders a client-side function stub for each route
    # emission
    -> std.template.render
  codegen.emit_migration
    fn (spec: app_spec) -> result[string, string]
    + renders a schema migration covering all entities
    # emission
    -> std.template.render
  codegen.project_plan
    fn (spec: app_spec) -> list[generated_file]
    + lists the relative paths and contents that would be emitted
    # planning
  codegen.write_project
    fn (plan: list[generated_file], out_dir: string) -> result[void, string]
    + writes every planned file under the output directory
    - returns error on filesystem failure
    # output
    -> std.fs.make_dir_all
    -> std.fs.write_all
