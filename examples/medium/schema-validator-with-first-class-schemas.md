# Requirement: "a schema validation library where schemas are first-class values that validate arbitrary structured input"

Schemas are built compositionally and then applied to a dynamic value tree. The project layer is the schema algebra plus a single validate entry point.

std
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses JSON into a dynamic value tree
      - returns error on invalid JSON
      # parsing

schema
  schema.string_schema
    fn (min_len: i32, max_len: i32) -> schema_node
    + builds a string constraint node
    # builders
  schema.number_schema
    fn (min: f64, max: f64) -> schema_node
    + builds a numeric range node
    # builders
  schema.bool_schema
    fn () -> schema_node
    + builds a boolean node
    # builders
  schema.object_schema
    fn (fields: map[string, schema_node], required: list[string]) -> schema_node
    + builds an object node with required keys
    # builders
  schema.array_schema
    fn (item: schema_node, min_items: i32, max_items: i32) -> schema_node
    + builds an array node with per-item constraints
    # builders
  schema.optional
    fn (inner: schema_node) -> schema_node
    + wraps a node so null is accepted
    # builders
  schema.validate
    fn (root: schema_node, value: json_value) -> result[void, list[string]]
    + returns ok when the value satisfies the schema
    - returns a list of path-scoped error messages when it does not
    # validation
  schema.validate_raw
    fn (root: schema_node, raw_json: string) -> result[void, list[string]]
    + parses raw JSON then validates against the schema
    - returns a parse error as a single entry
    # validation
    -> std.json.parse
    -> schema.validate
