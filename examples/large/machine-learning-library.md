# Requirement: "a machine learning library"

A general-purpose ML library exposing classical algorithms behind a fit/predict interface. The project layer names a handful of estimators; real numerics live in std.

std
  std.math
    std.math.vector_dot
      fn (a: list[f64], b: list[f64]) -> f64
      + returns the inner product of two equal-length vectors
      - returns 0.0 when either is empty
      # linear_algebra
    std.math.euclidean_distance
      fn (a: list[f64], b: list[f64]) -> f64
      + returns sqrt(sum((a_i - b_i)^2))
      # linear_algebra
    std.math.mean
      fn (xs: list[f64]) -> f64
      + returns arithmetic mean
      - returns 0.0 on empty input
      # statistics
    std.math.standard_deviation
      fn (xs: list[f64]) -> f64
      + returns sample standard deviation
      # statistics
  std.random
    std.random.shuffle_indices
      fn (n: i32, seed: i64) -> list[i32]
      + returns a deterministic shuffled permutation of [0, n)
      # randomness
  std.linalg
    std.linalg.matvec
      fn (m: list[list[f64]], v: list[f64]) -> list[f64]
      + multiplies a matrix by a column vector
      - returns empty list on dimension mismatch
      # linear_algebra

ml
  ml.train_test_split
    fn (x: list[list[f64]], y: list[f64], test_ratio: f64, seed: i64) -> tuple[list[list[f64]], list[f64], list[list[f64]], list[f64]]
    + splits samples into train and test portions by ratio
    - returns empty test set when ratio is 0
    # dataset_split
    -> std.random.shuffle_indices
  ml.standardize
    fn (x: list[list[f64]]) -> tuple[list[list[f64]], list[f64], list[f64]]
    + returns z-scored features along with per-column mean and stddev
    # preprocessing
    -> std.math.mean
    -> std.math.standard_deviation
  ml.linear_regression_fit
    fn (x: list[list[f64]], y: list[f64]) -> result[list[f64], string]
    + returns coefficients via normal equations
    - returns error when x is empty or rank deficient
    # model_fitting
    -> std.linalg.matvec
  ml.linear_regression_predict
    fn (coefs: list[f64], x: list[list[f64]]) -> list[f64]
    + returns predictions for each sample
    # model_inference
    -> std.math.vector_dot
  ml.knn_fit
    fn (x: list[list[f64]], y: list[f64], k: i32) -> knn_model
    + stores training set and k
    - returns model with empty data when x is empty
    # model_fitting
  ml.knn_predict
    fn (model: knn_model, x: list[list[f64]]) -> list[f64]
    + returns majority label of k nearest neighbors for each sample
    # model_inference
    -> std.math.euclidean_distance
  ml.accuracy_score
    fn (y_true: list[f64], y_pred: list[f64]) -> f64
    + returns fraction of exact matches
    - returns 0.0 when lengths differ
    # metrics
  ml.mean_squared_error
    fn (y_true: list[f64], y_pred: list[f64]) -> f64
    + returns mean of squared residuals
    # metrics
