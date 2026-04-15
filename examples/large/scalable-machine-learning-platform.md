# Requirement: "a scalable machine learning platform"

Core training/prediction loop for linear models with dataset handling and evaluation. Scope kept to one model family on purpose.

std
  std.math
    std.math.sigmoid
      fn (x: f64) -> f64
      + returns 1 / (1 + exp(-x))
      # math
    std.math.dot_product
      fn (a: list[f64], b: list[f64]) -> f64
      + returns sum of pairwise products
      - returns 0 when lengths mismatch
      # math
  std.rand
    std.rand.shuffle
      fn (n: i32, seed: i64) -> list[i32]
      + returns a permutation of 0..n under the given seed
      # randomness

ml_platform
  ml_platform.load_csv
    fn (path: string, label_column: i32) -> result[dataset, string]
    + parses a csv into numeric features and labels
    - returns error when a row has fewer columns than label_column
    # data_loading
  ml_platform.train_test_split
    fn (ds: dataset, test_fraction: f64, seed: i64) -> tuple[dataset, dataset]
    + returns (train, test) using the given fraction for test
    ? split is deterministic per seed
    # data_loading
    -> std.rand.shuffle
  ml_platform.standardize
    fn (ds: dataset) -> dataset
    + returns a dataset where each feature is zero-mean, unit-variance
    # preprocessing
  ml_platform.new_logistic_model
    fn (num_features: i32) -> logistic_model
    + returns a model with zeroed weights and bias
    # construction
  ml_platform.predict_proba
    fn (model: logistic_model, features: list[f64]) -> f64
    + returns sigmoid(w . x + b)
    # inference
    -> std.math.dot_product
    -> std.math.sigmoid
  ml_platform.train_sgd
    fn (model: logistic_model, ds: dataset, lr: f64, epochs: i32) -> logistic_model
    + runs stochastic gradient descent for the given epochs
    ? uses log-loss; gradient is (pred - label) * x
    # training
    -> std.rand.shuffle
  ml_platform.evaluate_accuracy
    fn (model: logistic_model, ds: dataset) -> f64
    + returns fraction of examples whose predicted class matches the label
    - returns 0 when the dataset is empty
    # evaluation
  ml_platform.save_model
    fn (model: logistic_model, path: string) -> result[void, string]
    + persists weights and bias to disk
    # persistence
  ml_platform.load_model
    fn (path: string) -> result[logistic_model, string]
    + rehydrates a model from disk
    - returns error on missing or corrupt file
    # persistence
