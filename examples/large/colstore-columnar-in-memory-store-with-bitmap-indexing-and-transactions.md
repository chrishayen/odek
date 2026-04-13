# Requirement: "a columnar embeddable in-memory store with bitmap indexing and transactions"

Column-oriented storage where each column is a typed vector and bitmap indexes mark which rows satisfy predicates. Transactions provide snapshot reads and atomic commits.

std
  std.collections
    std.collections.bitmap_new
      @ (capacity: u64) -> bitmap_state
      + returns a zeroed bitmap with the given bit capacity
      # collections
    std.collections.bitmap_set
      @ (bm: bitmap_state, index: u64, value: bool) -> bitmap_state
      + sets or clears a single bit
      # collections
    std.collections.bitmap_and
      @ (a: bitmap_state, b: bitmap_state) -> bitmap_state
      + returns the bitwise AND of two bitmaps
      # collections
    std.collections.bitmap_or
      @ (a: bitmap_state, b: bitmap_state) -> bitmap_state
      + returns the bitwise OR of two bitmaps
      # collections
    std.collections.bitmap_iter_set
      @ (bm: bitmap_state) -> list[u64]
      + returns the indices of every set bit in ascending order
      # collections

colstore
  colstore.open
    @ () -> store_state
    + returns an empty store with no columns and no rows
    # construction
  colstore.add_column
    @ (store: store_state, name: string, dtype: string) -> result[store_state, string]
    + registers a column with a type tag like "i64", "f64", or "string"
    - returns error on duplicate column name
    - returns error on unknown type
    # schema
  colstore.begin
    @ (store: store_state) -> txn_state
    + returns a read-write transaction over the current snapshot
    ? readers see a consistent snapshot until commit
    # transactions
  colstore.insert_row
    @ (txn: txn_state, values: map[string, string]) -> result[txn_state, string]
    + appends a row across all columns
    - returns error when a required column is missing from values
    # mutation
  colstore.commit
    @ (store: store_state, txn: txn_state) -> result[store_state, string]
    + atomically applies the transaction's writes
    - returns error on write-write conflict with another committed transaction
    # transactions
  colstore.rollback
    @ (txn: txn_state) -> void
    + discards the transaction with no side effects
    # transactions
  colstore.build_index
    @ (store: store_state, column: string) -> result[store_state, string]
    + builds a bitmap index mapping distinct values to row-bitmaps
    - returns error when the column is unknown
    # indexing
    -> std.collections.bitmap_new
    -> std.collections.bitmap_set
  colstore.query_eq
    @ (store: store_state, column: string, value: string) -> result[list[u64], string]
    + returns row ids whose column equals the given value, using the bitmap index when present
    - returns error when the column is unknown
    # query
    -> std.collections.bitmap_iter_set
  colstore.query_and
    @ (store: store_state, predicates: list[tuple[string, string]]) -> result[list[u64], string]
    + intersects bitmap results from multiple column-equality predicates
    - returns error when any column is unknown
    # query
    -> std.collections.bitmap_and
    -> std.collections.bitmap_iter_set
  colstore.project
    @ (store: store_state, row_ids: list[u64], columns: list[string]) -> result[list[map[string, string]], string]
    + returns the requested columns for the given row ids
    - returns error when a requested column is unknown
    # projection
