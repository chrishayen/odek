# Requirement: "a form-data binder that fills typed fields from a string map"

Binds a parsed form (string keys to string values) into a target schema of typed fields.

std: (all units exist)

bind
  bind.new_schema
    fn () -> schema_state
    + creates an empty target schema
    # construction
  bind.field_string
    fn (schema: schema_state, name: string, required: bool) -> schema_state
    + declares a string field by form key
    # schema
  bind.field_int
    fn (schema: schema_state, name: string, required: bool) -> schema_state
    + declares an integer field that will be parsed from its string value
    # schema
  bind.field_bool
    fn (schema: schema_state, name: string, required: bool) -> schema_state
    + declares a bool field (accepts "true"/"false"/"1"/"0")
    # schema
  bind.bind_form
    fn (schema: schema_state, form: map[string, string]) -> result[bound_values, list[bind_error]]
    + returns populated values when every required field is present and parses cleanly
    - returns errors listing each missing required field
    - returns errors listing each field that failed to parse
    # binding
  bind.get_string
    fn (values: bound_values, name: string) -> optional[string]
    + returns the bound string for the field if it was set
    # access
  bind.get_int
    fn (values: bound_values, name: string) -> optional[i64]
    + returns the bound integer for the field if it was set
    # access
  bind.get_bool
    fn (values: bound_values, name: string) -> optional[bool]
    + returns the bound bool for the field if it was set
    # access
