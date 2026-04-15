# Requirement: "a gradient boosting library over decision trees"

A forest of shallow regression trees fit on residuals; std supplies the numeric primitives the tree fitter repeatedly calls.

std
  std.math
    std.math.sum_f64
      fn (values: list[f64]) -> f64
      + returns the sum of every element
      + returns 0.0 for an empty list
      # math
    std.math.mean_f64
      fn (values: list[f64]) -> result[f64, string]
      + returns the arithmetic mean
      - returns error on empty input
      # math
    std.math.variance_f64
      fn (values: list[f64]) -> result[f64, string]
      + returns the population variance
      - returns error on empty input
      # math
    std.math.argsort_f64
      fn (values: list[f64]) -> list[i32]
      + returns indices that would sort the list in ascending order
      # math

gbt
  gbt.new
    fn (learning_rate: f64, max_depth: i32) -> gbt_state
    + creates an empty model with hyperparameters
    # construction
  gbt.fit_stump
    fn (features: list[list[f64]], residuals: list[f64]) -> result[tree_node, string]
    + fits a single decision stump that best reduces squared error on residuals
    - returns error when features and residuals have different lengths
    # tree_fitting
    -> std.math.argsort_f64
    -> std.math.variance_f64
    -> std.math.sum_f64
  gbt.fit_tree
    fn (features: list[list[f64]], residuals: list[f64], max_depth: i32) -> result[tree_node, string]
    + grows a regression tree up to max_depth by splitting on squared error
    - returns error on shape mismatch
    # tree_fitting
    -> std.math.variance_f64
  gbt.predict_tree
    fn (tree: tree_node, row: list[f64]) -> f64
    + returns the leaf value for a single feature row
    # prediction
  gbt.add_tree
    fn (state: gbt_state, tree: tree_node) -> gbt_state
    + appends a tree scaled by the learning rate to the ensemble
    # training
  gbt.fit
    fn (state: gbt_state, features: list[list[f64]], targets: list[f64], n_rounds: i32) -> result[gbt_state, string]
    + fits n_rounds of boosted trees against the targets
    - returns error on shape mismatch between features and targets
    # training
    -> std.math.mean_f64
  gbt.predict
    fn (state: gbt_state, row: list[f64]) -> f64
    + returns the sum of every tree's prediction plus the base value
    # prediction
  gbt.predict_batch
    fn (state: gbt_state, features: list[list[f64]]) -> list[f64]
    + returns predictions for every row
    # prediction
  gbt.feature_importance
    fn (state: gbt_state) -> list[f64]
    + returns per-feature importance as accumulated error reduction
    # inspection
    -> std.math.sum_f64
