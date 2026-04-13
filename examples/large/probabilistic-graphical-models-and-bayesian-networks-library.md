# Requirement: "a library for probabilistic graphical models and Bayesian networks"

Build discrete Bayesian networks, attach conditional probability tables, then run inference and structure/parameter learning.

std
  std.math
    std.math.log_f64
      @ (x: f64) -> f64
      + returns the natural logarithm of x
      - returns negative infinity for x <= 0
      # math
    std.math.exp_f64
      @ (x: f64) -> f64
      + returns e raised to x
      # math
  std.random
    std.random.new_rng
      @ (seed: u64) -> rng_state
      + returns a deterministic PRNG seeded with the given value
      # randomness
    std.random.uniform_f64
      @ (rng: rng_state) -> tuple[f64, rng_state]
      + returns a uniform sample in [0, 1) and the advanced state
      # randomness

bayes
  bayes.new_network
    @ () -> network_state
    + creates an empty Bayesian network
    # construction
  bayes.add_variable
    @ (state: network_state, name: string, values: list[string]) -> network_state
    + adds a discrete variable with the given possible values
    - panics when the name is already present
    # model
  bayes.add_edge
    @ (state: network_state, parent: string, child: string) -> result[network_state, string]
    + adds a directed edge between two existing variables
    - returns error when the edge would create a cycle
    # model
  bayes.set_cpt
    @ (state: network_state, variable: string, table: cpt) -> result[network_state, string]
    + attaches a conditional probability table for the variable
    - returns error when the table shape does not match the variable's parents
    # model
  bayes.topological_order
    @ (state: network_state) -> list[string]
    + returns variables in a topological order
    # analysis
  bayes.joint_probability
    @ (state: network_state, assignment: map[string, string]) -> result[f64, string]
    + returns the joint probability of a full assignment
    - returns error when the assignment is incomplete
    # inference
    -> std.math.log_f64
    -> std.math.exp_f64
  bayes.variable_elimination
    @ (state: network_state, query: list[string], evidence: map[string, string]) -> result[map[list[string], f64], string]
    + returns the posterior distribution over the query variables given evidence
    - returns error when a query or evidence variable is unknown
    # inference
  bayes.gibbs_sample
    @ (state: network_state, rng: rng_state, evidence: map[string, string], iterations: i32) -> result[map[string, f64], string]
    + returns marginals estimated via Gibbs sampling with the given evidence
    - returns error when iterations is non-positive
    # inference
    -> std.random.uniform_f64
  bayes.learn_parameters_mle
    @ (state: network_state, dataset: list[map[string, string]]) -> result[network_state, string]
    + fits every CPT by maximum likelihood on the dataset
    - returns error when any row is missing a variable
    # learning
  bayes.learn_structure_hill_climb
    @ (variables: list[variable_spec], dataset: list[map[string, string]]) -> result[network_state, string]
    + returns a network found by greedy score-based hill climbing
    - returns error when dataset is empty
    # learning
