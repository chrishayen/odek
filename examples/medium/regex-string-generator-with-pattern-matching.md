# Requirement: "a generator that produces random strings matching a given regular expression"

Parses a small regex subset into an AST and generates a random string that would match it. Supports literals, character classes, alternation, concatenation, and bounded repetition.

std
  std.random
    std.random.int_range
      fn (rng: rng_state, low: i32, high: i32) -> tuple[i32, rng_state]
      + returns a uniformly drawn integer in [low, high] and the advanced RNG state
      # random
    std.random.new
      fn (seed: u64) -> rng_state
      + creates a deterministic RNG state from a seed
      # random

regexgen
  regexgen.parse
    fn (pattern: string) -> result[regex_node, string]
    + parses the pattern into an AST covering literals, ".", character classes, alternation, concatenation, "?", "*", "+", and "{m,n}"
    - returns error on unbalanced brackets or parentheses
    - returns error on unsupported escape sequences
    # parsing
  regexgen.bound_repeat
    fn (node: regex_node, star_cap: i32) -> regex_node
    + replaces unbounded "*" and "+" quantifiers with "{0,cap}" and "{1,cap}" so generation terminates
    # bounding
  regexgen.generate_node
    fn (node: regex_node, rng: rng_state) -> tuple[string, rng_state]
    + produces one random matching string for the given node and returns the advanced RNG
    ? alternation picks uniformly among branches; classes pick uniformly among characters
    # generation
    -> std.random.int_range
  regexgen.generate
    fn (pattern: string, seed: u64, star_cap: i32) -> result[string, string]
    + parses, bounds, and generates one matching string from the seed
    - returns error when parsing fails
    # pipeline
    -> regexgen.parse
    -> regexgen.bound_repeat
    -> regexgen.generate_node
    -> std.random.new
