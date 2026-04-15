# Requirement: "a toolkit for building and analyzing recommender systems"

User-item rating prediction with a baseline and matrix factorization, plus evaluation utilities.

std
  std.math
    std.math.sqrt
      fn (x: f64) -> f64
      + returns the square root of x
      - returns nan for x < 0
      # math
  std.rand
    std.rand.uniform_f64
      fn (low: f64, high: f64) -> f64
      + returns a uniformly distributed f64 in [low, high)
      # randomness

recsys
  recsys.load_ratings
    fn (triples: list[tuple[string, string, f32]]) -> dataset_handle
    + indexes (user, item, rating) triples and assigns internal integer ids
    # data
  recsys.train_test_split
    fn (data: dataset_handle, test_fraction: f64) -> tuple[dataset_handle, dataset_handle]
    + randomly partitions ratings into train and test sets
    ? assumes test_fraction is between 0 and 1
    # evaluation
    -> std.rand.uniform_f64
  recsys.fit_baseline
    fn (data: dataset_handle) -> model_handle
    + fits a global mean plus per-user and per-item bias terms
    # training
  recsys.fit_matrix_factorization
    fn (data: dataset_handle, factors: i32, epochs: i32, learning_rate: f64, regularization: f64) -> model_handle
    + fits biased matrix factorization by stochastic gradient descent
    # training
    -> std.rand.uniform_f64
  recsys.predict
    fn (model: model_handle, user: string, item: string) -> f64
    + returns the predicted rating for a user-item pair
    + falls back to the global mean when either id is unknown
    # inference
  recsys.top_n
    fn (model: model_handle, user: string, candidates: list[string], n: i32) -> list[tuple[string, f64]]
    + returns the n highest-scoring items for a user, sorted descending
    # recommendations
  recsys.rmse
    fn (model: model_handle, test: dataset_handle) -> f64
    + returns the root-mean-square error of predictions on the test set
    # metrics
    -> std.math.sqrt
  recsys.mae
    fn (model: model_handle, test: dataset_handle) -> f64
    + returns the mean absolute error of predictions on the test set
    # metrics
