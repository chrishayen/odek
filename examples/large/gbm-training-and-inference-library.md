# Requirement: "a gradient boosted decision tree training and inference library"

Builds an ensemble of regression trees by fitting each to the gradient of a loss function. Training supports histogram-based split finding.

std
  std.math
    std.math.log
      @ (x: f64) -> f64
      + returns the natural logarithm
      # math
    std.math.exp
      @ (x: f64) -> f64
      + returns e raised to x
      # math
  std.stats
    std.stats.quantile
      @ (values: list[f64], q: f64) -> f64
      + returns the q-quantile of values
      # statistics

gbm
  gbm.new_dataset
    @ (features: list[list[f64]], labels: list[f64]) -> dataset_state
    + packs feature rows and labels into a training set
    # data
  gbm.build_histograms
    @ (dataset: dataset_state, max_bins: i32) -> histogram_state
    + bins each feature into up to max_bins quantile-based bins
    # preprocessing
    -> std.stats.quantile
  gbm.squared_loss_gradient
    @ (predictions: list[f64], labels: list[f64]) -> list[f64]
    + returns per-row gradients for squared error
    # loss
  gbm.logistic_loss_gradient
    @ (predictions: list[f64], labels: list[f64]) -> list[f64]
    + returns per-row gradients for binary log loss
    # loss
    -> std.math.exp
  gbm.find_best_split
    @ (histograms: histogram_state, gradients: list[f64], node_indices: list[i32]) -> optional[split_info]
    + scans histogram bins and returns the split with the best gain
    - returns none when no split improves the objective
    # split_finding
  gbm.grow_tree
    @ (histograms: histogram_state, gradients: list[f64], max_depth: i32, min_leaf: i32) -> tree_node
    + grows one regression tree using histogram-based splits
    # tree_training
  gbm.train
    @ (dataset: dataset_state, num_rounds: i32, learning_rate: f64, max_depth: i32) -> model_state
    + fits num_rounds trees, updating predictions after each
    # training
  gbm.predict_row
    @ (model: model_state, features: list[f64]) -> f64
    + returns the ensemble prediction for one row
    # inference
  gbm.predict_batch
    @ (model: model_state, features: list[list[f64]]) -> list[f64]
    + returns predictions for many rows
    # inference
  gbm.feature_importance
    @ (model: model_state) -> list[f64]
    + returns total gain per feature across the ensemble
    # introspection
