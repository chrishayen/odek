# Requirement: "a relational test data generator that writes CSV files"

Defines tables with columns, generates synthetic rows obeying foreign key relationships, and writes one CSV per table.

std
  std.fs
    std.fs.write_all
      @ (path: string, data: string) -> result[void, string]
      + writes data to path
      # filesystem
  std.random
    std.random.new_seeded
      @ (seed: i64) -> rng_state
      + creates a deterministic RNG
      # randomness
    std.random.next_i64
      @ (rng: rng_state, lo: i64, hi: i64) -> tuple[i64, rng_state]
      + returns a value in [lo, hi]
      # randomness

data_gen
  data_gen.new_schema
    @ () -> schema_state
    + creates an empty schema
    # construction
  data_gen.add_table
    @ (schema: schema_state, name: string, columns: list[column_spec]) -> schema_state
    + registers a table with the given column specs
    ? each column spec includes a name, type, and generator rule
    # schema
  data_gen.add_foreign_key
    @ (schema: schema_state, from_table: string, from_column: string, to_table: string, to_column: string) -> schema_state
    + declares that from_table.from_column references to_table.to_column
    # relationships
  data_gen.generate
    @ (schema: schema_state, counts: map[string, i64], seed: i64) -> map[string, list[list[string]]]
    + generates the requested row count per table, honoring foreign keys
    ? parent tables are generated before their children
    # generation
    -> std.random.new_seeded
    -> std.random.next_i64
  data_gen.encode_csv
    @ (header: list[string], rows: list[list[string]]) -> string
    + returns a CSV string with the given header and rows
    # csv
  data_gen.write
    @ (schema: schema_state, data: map[string, list[list[string]]], out_dir: string) -> result[i32, string]
    + writes one CSV per table into out_dir and returns the count written
    # output
    -> std.fs.write_all
