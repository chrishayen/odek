# Requirement: "an evolutionary optimization library"

A genetic algorithm framework: an initial population, fitness evaluation, selection, crossover, mutation, and generation stepping. Problem-specific hooks are referenced by tag so the core stays generic.

std
  std.random
    std.random.uniform_f64
      @ () -> f64
      + returns a random value in [0, 1)
      # random
    std.random.int_range
      @ (low: i32, high: i32) -> i32
      + returns a random integer in [low, high)
      # random

evolve
  evolve.individual
    @ (genome: list[f64], fitness: f64) -> individual
    + builds an individual with its genome and last evaluated fitness
    # model
  evolve.population
    @ (individuals: list[individual]) -> population_state
    + wraps a slice of individuals as a population
    # model
  evolve.init_population
    @ (size: i32, genome_length: i32, low: f64, high: f64) -> population_state
    + creates a random population with genome values drawn uniformly in [low, high]
    -> std.random.uniform_f64
    # initialization
  evolve.evaluate
    @ (pop: population_state, fitness_tag: string) -> population_state
    + evaluates every individual using the named fitness function
    # evaluation
  evolve.select_tournament
    @ (pop: population_state, tournament_size: i32) -> individual
    + returns the best of k randomly sampled individuals
    ? picks individuals with replacement
    -> std.random.int_range
    # selection
  evolve.crossover_uniform
    @ (a: individual, b: individual) -> individual
    + produces a child by independently inheriting each gene from one parent
    -> std.random.uniform_f64
    # crossover
  evolve.mutate_gaussian
    @ (ind: individual, rate: f64, sigma: f64) -> individual
    + perturbs each gene with probability `rate` by gaussian noise of stddev `sigma`
    -> std.random.uniform_f64
    # mutation
  evolve.step_generation
    @ (pop: population_state, fitness_tag: string, tournament_size: i32, mutation_rate: f64, mutation_sigma: f64) -> population_state
    + produces the next generation via selection, crossover, mutation, and re-evaluation
    # evolution
  evolve.best
    @ (pop: population_state) -> optional[individual]
    + returns the individual with the highest fitness
    - returns none when the population is empty
    # query
