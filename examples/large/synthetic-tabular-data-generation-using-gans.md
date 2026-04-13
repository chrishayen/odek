# Requirement: "a library for synthetic tabular data generation using generative models"

Infers a schema from a sample table, fits a simple generative model per column, and samples new rows. The generative model is a toy mixture-of-categoricals plus gaussian for numeric columns — rich generators are out of scope for a reference decomposition.

std
  std.random
    std.random.new_seeded
      @ (seed: u64) -> rng_state
      + creates a deterministic PRNG
      # random
    std.random.next_f64
      @ (rng: rng_state) -> tuple[f64, rng_state]
      + returns a uniform value in [0, 1) and the advanced state
      # random
    std.random.next_i32_in_range
      @ (rng: rng_state, low: i32, high: i32) -> tuple[i32, rng_state]
      + returns an integer uniformly in [low, high)
      # random
  std.math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the square root; NaN for negative input
      # math
    std.math.ln
      @ (x: f64) -> f64
      + returns the natural logarithm
      # math
    std.math.erf_inv
      @ (x: f64) -> f64
      + returns the inverse error function for x in (-1, 1)
      # math

tabgen
  tabgen.infer_schema
    @ (rows: list[map[string, string]]) -> schema_state
    + detects each column as numeric or categorical
    ? a column is numeric when every non-empty value parses as f64
    # schema
  tabgen.fit
    @ (schema: schema_state, rows: list[map[string, string]]) -> model_state
    + trains a per-column model: mean/variance for numeric, frequency for categorical
    # training
  tabgen.sample_column_numeric
    @ (model: model_state, column: string, rng: rng_state) -> result[tuple[f64, rng_state], string]
    + draws a gaussian value using the column's mean and variance
    - returns error when the column is not numeric
    # sampling
    -> std.random.next_f64
    -> std.math.sqrt
    -> std.math.erf_inv
  tabgen.sample_column_categorical
    @ (model: model_state, column: string, rng: rng_state) -> result[tuple[string, rng_state], string]
    + draws a categorical value proportional to observed frequencies
    - returns error when the column is not categorical
    # sampling
    -> std.random.next_f64
  tabgen.sample_row
    @ (model: model_state, rng: rng_state) -> tuple[map[string, string], rng_state]
    + samples a single synthetic row across all columns
    # sampling
  tabgen.sample_rows
    @ (model: model_state, count: i32, seed: u64) -> list[map[string, string]]
    + samples count synthetic rows deterministically from the seed
    # sampling
    -> std.random.new_seeded
  tabgen.column_stats
    @ (model: model_state, column: string) -> result[column_summary, string]
    + returns the learned parameters for the column
    - returns error when the column is unknown
    # introspection
  tabgen.columns
    @ (schema: schema_state) -> list[string]
    + returns the column names in declaration order
    # introspection
  tabgen.column_kind
    @ (schema: schema_state, column: string) -> result[string, string]
    + returns "numeric" or "categorical"
    - returns error when the column is unknown
    # introspection
