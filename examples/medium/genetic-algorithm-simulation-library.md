# Requirement: "a genetic algorithm simulation library"

The project exposes generic evolution primitives parameterized by fitness; crossover and mutation are pluggable hooks the caller supplies as genome operations.

std
  std.random
    std.random.uniform_f64
      fn (seed_state: rng_state) -> tuple[f64, rng_state]
      + returns a uniform random value in [0.0, 1.0) and advances the state
      # randomness
    std.random.range_i32
      fn (seed_state: rng_state, low: i32, high: i32) -> tuple[i32, rng_state]
      + returns a random integer in [low, high)
      - returns low when low >= high
      # randomness
    std.random.new_rng
      fn (seed: i64) -> rng_state
      + constructs an rng state from a seed
      # randomness

genetic
  genetic.new_population
    fn (genomes: list[bytes], fitnesses: list[f64]) -> population
    + pairs genomes with their fitness values
    - returns an empty population when inputs have unequal length
    # construction
  genetic.select_tournament
    fn (pop: population, tournament_size: i32, rng: rng_state) -> tuple[bytes, rng_state]
    + returns the fittest genome among tournament_size random picks
    + advances the rng state
    # selection
    -> std.random.range_i32
  genetic.crossover_single_point
    fn (parent_a: bytes, parent_b: bytes, rng: rng_state) -> tuple[bytes, rng_state]
    + returns a child spliced at a random byte index
    - returns parent_a unchanged when either parent is empty
    # recombination
    -> std.random.range_i32
  genetic.mutate_bitflip
    fn (genome: bytes, rate: f64, rng: rng_state) -> tuple[bytes, rng_state]
    + flips each bit with probability rate
    # mutation
    -> std.random.uniform_f64
  genetic.evolve_generation
    fn (pop: population, next_fitness: list[f64], rng: rng_state) -> tuple[population, rng_state]
    + produces the next generation via tournament selection, crossover, and mutation
    ? caller supplies fitness scores for the new children via next_fitness
    # generation_step
  genetic.best_of
    fn (pop: population) -> optional[bytes]
    + returns the genome with the highest fitness
    - returns none when population is empty
    # query
