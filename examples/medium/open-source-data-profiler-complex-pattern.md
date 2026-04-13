# Requirement: "a data profiler that discovers functional dependencies, unique columns, and inclusion dependencies in tabular data"

Loads a dataset and runs a suite of profilers to surface structural patterns.

std
  std.csv
    std.csv.parse
      @ (raw: bytes) -> result[table, string]
      + returns a table with header and string-valued rows
      - returns error on malformed CSV
      # parsing

profiler
  profiler.load_table
    @ (raw: bytes) -> result[table, string]
    + returns a table ready for profiling
    # ingestion
    -> std.csv.parse
  profiler.column_stats
    @ (t: table, column: string) -> column_stats
    + returns count, distinct_count, nulls, min, and max for the column
    # descriptive_stats
  profiler.find_unique_columns
    @ (t: table) -> list[string]
    + returns columns whose non-null values are pairwise distinct
    # uniqueness
  profiler.find_functional_dependencies
    @ (t: table, max_lhs_size: i32) -> list[functional_dependency]
    + returns minimal dependencies A -> B where A determines B
    ? candidates are pruned when a subset already determines B
    # fd_discovery
  profiler.find_inclusion_dependencies
    @ (left: table, right: table) -> list[inclusion_dependency]
    + returns column pairs (a, b) where all values of a appear in b
    # ind_discovery
  profiler.find_duplicates
    @ (t: table, columns: list[string]) -> list[list[i32]]
    + returns row-index groups that share identical values on the given columns
    # duplicate_detection
  profiler.summary
    @ (t: table) -> profile_summary
    + returns the aggregated profile including stats, uniques, and dependencies
    # reporting
