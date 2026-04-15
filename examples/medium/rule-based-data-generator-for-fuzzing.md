# Requirement: "a rule-based data generator for fuzzing"

A tiny grammar-style generator. Rules are named productions that expand into literals, alternatives, sequences, or references. A seeded PRNG drives all choices so runs are reproducible.

std
  std.random
    std.random.new_seeded
      fn (seed: i64) -> rng_state
      + creates a deterministic PRNG
      # random
    std.random.next_in_range
      fn (state: rng_state, lo: i32, hi_exclusive: i32) -> tuple[i32, rng_state]
      + returns a uniform integer in [lo, hi_exclusive) and the advanced state
      # random

fuzzgen
  fuzzgen.new_grammar
    fn () -> grammar
    + creates an empty grammar with no rules defined
    # construction
  fuzzgen.define_literal
    fn (g: grammar, name: string, value: string) -> grammar
    + registers a rule that always expands to the given literal
    # rule
  fuzzgen.define_alt
    fn (g: grammar, name: string, choices: list[string]) -> grammar
    + registers a rule that picks uniformly from a list of rule references
    # rule
  fuzzgen.define_sequence
    fn (g: grammar, name: string, parts: list[string]) -> grammar
    + registers a rule that concatenates the expansions of each referenced rule
    # rule
  fuzzgen.expand
    fn (g: grammar, start: string, state: rng_state, max_depth: i32) -> result[tuple[string, rng_state], string]
    + expands the named rule, returning the generated string and advanced RNG state
    - returns error when a referenced rule is undefined
    - returns error when expansion exceeds max_depth
    # expansion
    -> std.random.next_in_range
