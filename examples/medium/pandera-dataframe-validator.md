# Requirement: "a data validation library for dataframes"

Declarative schemas for tabular data: per-column type and constraint rules, applied to a dataframe-shaped input. The project layer defines the schema vocabulary and the check loop.

std: (all units exist)

pandera
  pandera.column
    fn (name: string, dtype: string) -> column_spec
    + creates a column spec requiring the given name and dtype
    ? dtype is a string tag like "i64", "f64", or "string"
    # schema
  pandera.with_nullable
    fn (spec: column_spec, nullable: bool) -> column_spec
    + sets whether the column permits null values
    # schema
  pandera.with_range
    fn (spec: column_spec, min: f64, max: f64) -> column_spec
    + attaches an inclusive numeric range constraint
    # schema
  pandera.with_regex
    fn (spec: column_spec, pattern: string) -> column_spec
    + attaches a regex constraint applied to every cell
    # schema
  pandera.with_unique
    fn (spec: column_spec) -> column_spec
    + requires the column to contain no duplicates
    # schema
  pandera.schema
    fn (columns: list[column_spec]) -> schema_state
    + assembles a set of column specs into a schema
    # schema
  pandera.validate
    fn (schema: schema_state, frame: dataframe) -> result[void, list[validation_error]]
    + returns Ok when the frame conforms
    - returns a list of validation errors with column name and row index when it does not
    - reports a schema error when a declared column is missing
    # validation
