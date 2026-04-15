# Requirement: "a dataframe library with columnar storage, typed columns, and basic transforms"

A dataframe is a set of named columns of equal length. Transforms are pure and return new dataframes.

std: (all units exist)

dataframe
  dataframe.empty
    fn () -> frame_state
    + creates a dataframe with no columns and zero rows
    # construction
  dataframe.from_columns
    fn (columns: map[string, column_state]) -> result[frame_state, string]
    + creates a dataframe from a set of named columns
    - returns error when columns have unequal lengths
    # construction
  dataframe.row_count
    fn (state: frame_state) -> i64
    + returns the number of rows in the dataframe
    # introspection
  dataframe.column_names
    fn (state: frame_state) -> list[string]
    + returns the column names in insertion order
    # introspection
  dataframe.select
    fn (state: frame_state, columns: list[string]) -> result[frame_state, string]
    + returns a dataframe containing only the named columns
    - returns error when a requested name is missing
    # projection
  dataframe.drop
    fn (state: frame_state, columns: list[string]) -> frame_state
    + returns a dataframe without the named columns
    # projection
  dataframe.rename
    fn (state: frame_state, mapping: map[string, string]) -> frame_state
    + returns a dataframe with columns renamed per the mapping
    # projection
  dataframe.filter_mask
    fn (state: frame_state, mask: list[bool]) -> result[frame_state, string]
    + returns a dataframe keeping rows where mask is true
    - returns error when mask length does not equal row count
    # filtering
  dataframe.sort_by
    fn (state: frame_state, column: string, ascending: bool) -> result[frame_state, string]
    + returns a dataframe sorted by the given column
    - returns error when the column is missing
    # sorting
  dataframe.group_sum
    fn (state: frame_state, key: string, value: string) -> result[frame_state, string]
    + returns a two-column dataframe summing value grouped by key
    - returns error when either column is missing
    # aggregation
  dataframe.concat
    fn (left: frame_state, right: frame_state) -> result[frame_state, string]
    + stacks two dataframes vertically
    - returns error when column sets do not match
    # combination
  dataframe.join_inner
    fn (left: frame_state, right: frame_state, key: string) -> result[frame_state, string]
    + returns the inner join on the given key column
    - returns error when the key is missing from either side
    # combination
  dataframe.column
    fn (state: frame_state, name: string) -> result[column_state, string]
    + returns the column data by name
    - returns error when the column does not exist
    # introspection
