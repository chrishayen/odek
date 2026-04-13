# Requirement: "a compile-time object-relational mapper generator driven by struct descriptors"

Accepts a normalized struct-to-table mapping and emits typed CRUD query stubs along with parameterized SQL. The caller supplies descriptors; the library does the analysis and code shape.

std
  std.strings
    std.strings.to_snake_case
      @ (s: string) -> string
      + converts a PascalCase name to snake_case
      # text_processing
    std.strings.join
      @ (parts: list[string], sep: string) -> string
      + joins the parts with the separator
      # text_processing

orm
  orm.new_entity
    @ (type_name: string) -> entity_descriptor
    + creates an entity descriptor with an inferred table name
    ? inferred table name is the snake_case form of the type name
    # construction
    -> std.strings.to_snake_case
  orm.set_table
    @ (e: entity_descriptor, table: string) -> entity_descriptor
    + overrides the default table name
    # construction
  orm.add_column
    @ (e: entity_descriptor, field: string, column: string, sql_type: string, primary_key: bool) -> entity_descriptor
    + appends a column mapping with its SQL type and primary-key flag
    # construction
  orm.validate
    @ (e: entity_descriptor) -> result[void, string]
    + checks that the entity has exactly one primary key and at least one column
    - returns error when there is zero or more than one primary key
    # validation
  orm.emit_select_by_id
    @ (e: entity_descriptor) -> string
    + renders a parameterized SELECT query keyed on the primary key column
    # sql_emission
    -> std.strings.join
  orm.emit_insert
    @ (e: entity_descriptor) -> string
    + renders a parameterized INSERT with one placeholder per non-primary-key column
    # sql_emission
    -> std.strings.join
  orm.emit_update_by_id
    @ (e: entity_descriptor) -> string
    + renders a parameterized UPDATE that sets non-primary-key columns and filters by the primary key
    # sql_emission
    -> std.strings.join
  orm.emit_delete_by_id
    @ (e: entity_descriptor) -> string
    + renders a parameterized DELETE by primary key
    # sql_emission
  orm.emit_typed_stub
    @ (e: entity_descriptor, op: string) -> result[string, string]
    + renders a typed function stub (find_by_id, insert, update, delete) whose parameters match the column field types
    - returns error when op is not one of the four supported operations
    # code_emission
  orm.generate
    @ (entities: list[entity_descriptor]) -> result[string, string]
    + validates each entity, then emits all four CRUD stubs and their SQL for each
    - returns error on the first invalid entity
    # pipeline
    -> orm.validate
    -> orm.emit_select_by_id
    -> orm.emit_insert
    -> orm.emit_update_by_id
    -> orm.emit_delete_by_id
    -> orm.emit_typed_stub
