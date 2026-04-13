# Requirement: "a dataframe structure with column-oriented operations"

A minimal dataframe: typed columns, row access, filter, select, and aggregate.

std: (all units exist)

dataframe
  dataframe.new
    @ (columns: list[column_spec]) -> dataframe
    + creates an empty dataframe with the given named, typed columns
    # construction
  dataframe.append_row
    @ (df: dataframe, row: list[cell]) -> result[dataframe, string]
    + returns a new dataframe with the row appended
    - returns error when row length does not match column count
    - returns error when any cell type does not match its column
    # mutation
  dataframe.select
    @ (df: dataframe, names: list[string]) -> result[dataframe, string]
    + returns a dataframe containing only the named columns
    - returns error when a requested column is missing
    # projection
  dataframe.filter
    @ (df: dataframe, predicate: row_predicate) -> dataframe
    + returns rows for which the predicate is true
    # filtering
  dataframe.group_by
    @ (df: dataframe, key_col: string) -> result[map[string, dataframe], string]
    + returns a map from key value to the subset of rows with that key
    - returns error when the key column does not exist
    # grouping
  dataframe.aggregate_sum
    @ (df: dataframe, col: string) -> result[f64, string]
    + returns the sum of a numeric column
    - returns error when the column is non-numeric
    # aggregation
