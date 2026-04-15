# Requirement: "a framework for grammar-driven fuzzing paired with delta debugging of failing inputs"

Generates random inputs from a grammar, runs them through a user-supplied oracle, and shrinks any failing input by minimizing against the same grammar.

std
  std.random
    std.random.new
      fn (seed: u64) -> rng_state
      + creates a deterministic RNG state from a seed
      # random
    std.random.int_range
      fn (rng: rng_state, low: i32, high: i32) -> tuple[i32, rng_state]
      + returns a uniformly drawn integer in [low, high] and the advanced RNG
      # random

tavor
  tavor.new_grammar
    fn () -> grammar_state
    + creates an empty grammar
    # construction
  tavor.add_rule
    fn (g: grammar_state, name: string, productions: list[list[string]]) -> grammar_state
    + adds a non-terminal with one or more production alternatives of tokens
    ? tokens prefixed with "$" refer to other rules; anything else is a literal
    # construction
  tavor.validate_grammar
    fn (g: grammar_state, start: string) -> result[void, string]
    + checks that the start rule exists and every referenced rule is defined
    - returns error on the first undefined reference
    # validation
  tavor.generate
    fn (g: grammar_state, start: string, rng: rng_state, depth_cap: i32) -> tuple[string, rng_state]
    + expands the start rule, picking one alternative per non-terminal, until no non-terminals remain or depth_cap is exceeded
    ? at the cap the expander prefers shorter alternatives to terminate
    # generation
    -> std.random.int_range
  tavor.fuzz
    fn (g: grammar_state, start: string, seed: u64, iterations: i32, oracle: oracle_fn) -> optional[string]
    + generates inputs until the oracle rejects one, then returns that failing input
    - returns none after iterations generations with no rejection
    # fuzzing
    -> std.random.new
    -> tavor.generate
  tavor.tree_of
    fn (g: grammar_state, start: string, input: string) -> optional[derivation_tree]
    + recovers a derivation tree for the input under the grammar
    - returns none when the input cannot be derived
    # parsing
  tavor.shrink_tree
    fn (t: derivation_tree, oracle: oracle_fn) -> derivation_tree
    + greedily replaces sub-trees with smaller alternatives while the oracle continues to reject
    # delta_debugging
  tavor.minimize
    fn (g: grammar_state, start: string, failing: string, oracle: oracle_fn) -> result[string, string]
    + parses the failing input against the grammar and returns the smallest still-failing rendering
    - returns error when the input does not parse
    # delta_debugging
    -> tavor.tree_of
    -> tavor.shrink_tree
