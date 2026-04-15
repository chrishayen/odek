# Requirement: "a genetic algorithm library"

Fixed-length bitstring individuals. The caller supplies a fitness function handle; the library handles selection, crossover, mutation, and generations.

std
  std.random
    std.random.u32
      fn () -> u32
      + returns a uniform random 32-bit unsigned integer
      # randomness
    std.random.f64_unit
      fn () -> f64
      + returns a uniform random double in [0.0, 1.0)
      # randomness

genetic
  genetic.new_population
    fn (size: i32, genome_length: i32) -> population_state
    + creates a population of `size` individuals with random bitstring genomes of the given length
    # construction
    -> std.random.u32
  genetic.evaluate
    fn (state: population_state, fitness_fn: string) -> population_state
    + invokes the caller-supplied fitness handle for each individual and stores the scores
    ? fitness_fn is an opaque handle the runtime resolves
    # evaluation
  genetic.select_tournament
    fn (state: population_state, tournament_size: i32) -> i32
    + picks `tournament_size` random individuals and returns the index of the fittest
    # selection
    -> std.random.u32
  genetic.crossover_single_point
    fn (parent_a: bytes, parent_b: bytes) -> bytes
    + returns a child genome by splicing the two parents at a random point
    # crossover
    -> std.random.u32
  genetic.mutate
    fn (genome: bytes, rate: f64) -> bytes
    + flips each bit with probability `rate`
    # mutation
    -> std.random.f64_unit
  genetic.step_generation
    fn (state: population_state, fitness_fn: string, tournament_size: i32, mutation_rate: f64) -> population_state
    + produces the next generation via tournament selection, crossover, and mutation
    + preserves population size
    # generation_step
  genetic.best
    fn (state: population_state) -> tuple[bytes, f64]
    + returns the fittest genome and its score
    # introspection
