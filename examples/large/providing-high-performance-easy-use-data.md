# Requirement: "a tabular data structures and analysis library"

Columnar tables with typed columns, selection, aggregation, and join primitives.

std
  std.math
    std.math.sum_f64
      @ (xs: list[f64]) -> f64
      + returns the sum of the values
      + returns 0 for an empty list
      # math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the non-negative square root
      # math
  std.sort
    std.sort.indices_by_f64
      @ (xs: list[f64], ascending: bool) -> list[i64]
      + returns the stable sort permutation of indices
      # sorting

dataframe
  dataframe.new
    @ (column_names: list[string]) -> dataframe_state
    + builds an empty frame with the given column schema
    # construction
  dataframe.set_column_f64
    @ (df: dataframe_state, name: string, values: list[f64]) -> result[dataframe_state, string]
    + attaches the values to the named column
    - returns error when the name is not in the schema
    - returns error when the length does not match existing columns
    # mutation
  dataframe.set_column_string
    @ (df: dataframe_state, name: string, values: list[string]) -> result[dataframe_state, string]
    + attaches the string values to the named column
    - returns error when the name is not in the schema
    # mutation
  dataframe.row_count
    @ (df: dataframe_state) -> i64
    + returns the number of rows in the frame
    # introspection
  dataframe.select_columns
    @ (df: dataframe_state, names: list[string]) -> result[dataframe_state, string]
    + returns a new frame with only the requested columns in order
    - returns error when any name is missing
    # projection
  dataframe.filter_rows
    @ (df: dataframe_state, predicate: fn(row_index: i64) -> bool) -> dataframe_state
    + returns a frame containing only the rows where the predicate holds
    # selection
  dataframe.sort_by
    @ (df: dataframe_state, column: string, ascending: bool) -> result[dataframe_state, string]
    + returns the frame with rows reordered by the column
    - returns error when the column does not exist
    # ordering
    -> std.sort.indices_by_f64
  dataframe.group_by_sum
    @ (df: dataframe_state, key: string, value: string) -> result[dataframe_state, string]
    + returns a frame with one row per unique key and the sum of the value column
    - returns error when either column is missing
    # aggregation
    -> std.math.sum_f64
  dataframe.mean
    @ (df: dataframe_state, column: string) -> result[f64, string]
    + returns the arithmetic mean of the numeric column
    - returns error when the column is non-numeric or missing
    - returns error when the frame is empty
    # aggregation
    -> std.math.sum_f64
  dataframe.std_dev
    @ (df: dataframe_state, column: string) -> result[f64, string]
    + returns the sample standard deviation of the numeric column
    - returns error when the frame has fewer than two rows
    # aggregation
    -> std.math.sqrt
  dataframe.inner_join
    @ (left: dataframe_state, right: dataframe_state, key: string) -> result[dataframe_state, string]
    + returns a frame of rows matching on the key column
    - returns error when either side lacks the key column
    # join
