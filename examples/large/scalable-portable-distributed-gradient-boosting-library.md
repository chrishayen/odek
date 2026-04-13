# Requirement: "a scalable, portable, distributed gradient boosting library"

Gradient boosted decision trees with histogram-based splits and distributed parameter aggregation.

std
  std.math
    std.math.exp
      @ (x: f64) -> f64
      + returns e raised to x
      # math
    std.math.log
      @ (x: f64) -> f64
      + returns the natural logarithm of x
      - returns nan for x <= 0
      # math
  std.net
    std.net.allreduce
      @ (local: list[f64]) -> result[list[f64], string]
      + sums the vector element-wise across all participating workers
      - returns error when the collective cannot complete
      # collectives

gboost
  gboost.new_dataset
    @ (features: list[list[f32]], labels: list[f32]) -> result[dataset_handle, string]
    + creates a dataset from dense rows and matching labels
    - returns error when row count does not match label count
    # data
  gboost.build_histograms
    @ (data: dataset_handle, max_bins: i32) -> dataset_handle
    + bins each feature column into at most max_bins discrete buckets
    ? binning is done once up front so splits operate on bin indices
    # preprocessing
  gboost.sigmoid
    @ (x: f64) -> f64
    + returns 1 / (1 + exp(-x))
    # activations
    -> std.math.exp
  gboost.logloss_gradient
    @ (pred: f64, label: f32) -> tuple[f64, f64]
    + returns (gradient, hessian) of binary log loss at the current prediction
    # objectives
    -> std.math.log
  gboost.find_best_split
    @ (data: dataset_handle, gradients: list[f64], hessians: list[f64], feature_idx: i32) -> split_info
    + scans histogram buckets to find the split that maximizes gain
    # tree_learning
  gboost.grow_tree
    @ (data: dataset_handle, gradients: list[f64], hessians: list[f64], max_depth: i32) -> tree_handle
    + grows one regression tree by recursively choosing best splits until max_depth
    # tree_learning
  gboost.predict_tree
    @ (tree: tree_handle, row: list[f32]) -> f64
    + returns the leaf value for a row by walking splits from the root
    # inference
  gboost.train
    @ (data: dataset_handle, num_rounds: i32, learning_rate: f64, max_depth: i32) -> result[model_handle, string]
    + fits num_rounds boosted trees by repeatedly computing gradients and growing trees
    - returns error when num_rounds < 1
    # training
  gboost.distributed_aggregate
    @ (local_histograms: list[f64]) -> result[list[f64], string]
    + sums per-worker histograms so all workers see the global split statistics
    # distribution
    -> std.net.allreduce
  gboost.predict
    @ (model: model_handle, row: list[f32]) -> f64
    + returns the summed leaf contributions across all trees
    # inference
  gboost.save_model
    @ (model: model_handle) -> bytes
    + serializes a model to a portable byte buffer
    # serialization
  gboost.load_model
    @ (buf: bytes) -> result[model_handle, string]
    + reconstructs a model from a portable byte buffer
    - returns error on a malformed buffer
    # serialization
