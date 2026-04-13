# Requirement: "a Bradley-Terry model for pairwise comparisons"

Fits per-item strength parameters from pairwise win/loss counts using iterative maximum-likelihood updates.

std: (all units exist)

bradley_terry
  bradley_terry.new
    @ (item_ids: list[string]) -> bt_state
    + creates a model with one strength parameter per item, all initialized to 1.0
    # construction
  bradley_terry.record_match
    @ (state: bt_state, winner: string, loser: string) -> bt_state
    + increments the win count for (winner, loser)
    # data_ingest
  bradley_terry.fit
    @ (state: bt_state, max_iterations: i32, tolerance: f64) -> bt_state
    + runs MM iterative updates until strengths change by less than tolerance or the iteration cap is reached
    ? strengths are normalized so their product equals 1 after each iteration
    # fitting
  bradley_terry.win_probability
    @ (state: bt_state, a: string, b: string) -> result[f64, string]
    + returns the probability that a beats b under the fitted model
    - returns error when either id is not in the model
    # prediction
