# Requirement: "a feature engineering and selection library with a fit/transform API"

A small set of transformers with a uniform fit/transform shape, plus one selector. Each transformer carries learned parameters; transform is a pure application of those parameters.

std: (all units exist)

features
  features.fit_standard_scaler
    fn (columns: list[list[f64]]) -> result[scaler_params, string]
    + returns per-column mean and standard deviation
    - returns error when columns is empty
    - returns error when rows have inconsistent lengths
    # fitting
  features.transform_standard_scaler
    fn (params: scaler_params, columns: list[list[f64]]) -> result[list[list[f64]], string]
    + returns columns with (x - mean) / std applied
    - returns error when column count does not match the fitted params
    # transformation
  features.fit_min_max_scaler
    fn (columns: list[list[f64]]) -> result[min_max_params, string]
    + returns per-column min and max
    - returns error when columns is empty
    # fitting
  features.transform_min_max_scaler
    fn (params: min_max_params, columns: list[list[f64]]) -> result[list[list[f64]], string]
    + returns columns rescaled to [0, 1]
    + returns 0.5 for a column whose min equals its max
    # transformation
  features.fit_one_hot_encoder
    fn (column: list[string]) -> one_hot_params
    + returns the sorted set of unique category labels
    # fitting
  features.transform_one_hot_encoder
    fn (params: one_hot_params, column: list[string]) -> list[list[f64]]
    + returns one row per input with a 1.0 in the matching category column
    + returns all zeros for labels unseen during fit
    # transformation
  features.select_by_variance_threshold
    fn (columns: list[list[f64]], threshold: f64) -> result[list[i32], string]
    + returns indices of columns whose variance is strictly greater than threshold
    - returns error when columns is empty
    # selection
