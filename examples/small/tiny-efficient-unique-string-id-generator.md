# Requirement: "a tiny and efficient unique string id generator"

Generates short random IDs from a fixed alphabet. Randomness is a thin std primitive so tests can inject a deterministic source.

std
  std.random
    std.random.fill_bytes
      fn (n: i32) -> bytes
      + returns n cryptographically random bytes
      # randomness

id_generator
  id_generator.generate
    fn (size: i32) -> string
    + returns a string of the given length drawn from the default url-safe alphabet
    + successive calls return distinct values with overwhelming probability
    - returns empty string when size <= 0
    # id_generation
    -> std.random.fill_bytes
  id_generator.generate_with_alphabet
    fn (alphabet: string, size: i32) -> string
    + returns a string of the given length drawn uniformly from the alphabet
    - returns empty string when alphabet is empty or size <= 0
    ? uses rejection sampling to avoid modulo bias
    # id_generation
    -> std.random.fill_bytes
