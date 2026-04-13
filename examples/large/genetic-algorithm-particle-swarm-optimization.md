# Requirement: "a genetic algorithm and particle swarm optimization library"

Two optimizers over fixed-dimensional real-valued vectors, sharing a common fitness-evaluation shape. Fitness functions are opaque caller-supplied handles.

std
  std.random
    std.random.f64_unit
      @ () -> f64
      + returns a uniform random double in [0.0, 1.0)
      # randomness
    std.random.f64_range
      @ (lo: f64, hi: f64) -> f64
      + returns a uniform random double in [lo, hi)
      # randomness
    std.random.u32
      @ () -> u32
      + returns a uniform random 32-bit unsigned integer
      # randomness

evoli
  evoli.new_ga
    @ (pop_size: i32, dims: i32, lo: f64, hi: f64) -> ga_state
    + creates a GA population with real-valued vectors of length `dims` uniformly in [lo, hi]
    # ga_construction
    -> std.random.f64_range
  evoli.ga_evaluate
    @ (state: ga_state, fitness_fn: string) -> ga_state
    + evaluates each individual with the caller's fitness handle and stores scores
    # evaluation
  evoli.ga_select
    @ (state: ga_state, tournament_size: i32) -> i32
    + tournament selection returning the index of the winner
    # selection
    -> std.random.u32
  evoli.ga_crossover
    @ (a: list[f64], b: list[f64], blend: f64) -> list[f64]
    + returns a child vector as blend*a + (1-blend)*b
    ? real-valued linear crossover; `blend` is typically sampled per-call
    # crossover
  evoli.ga_mutate
    @ (vec: list[f64], rate: f64, sigma: f64) -> list[f64]
    + with probability `rate` per coordinate, adds gaussian noise with std `sigma`
    # mutation
    -> std.random.f64_unit
  evoli.ga_step
    @ (state: ga_state, fitness_fn: string, tournament_size: i32, mutation_rate: f64, sigma: f64) -> ga_state
    + produces the next GA generation
    # generation_step
  evoli.new_pso
    @ (swarm_size: i32, dims: i32, lo: f64, hi: f64) -> pso_state
    + creates a PSO swarm with random positions and zero velocities
    # pso_construction
    -> std.random.f64_range
  evoli.pso_evaluate
    @ (state: pso_state, fitness_fn: string) -> pso_state
    + evaluates each particle and updates personal bests
    # evaluation
  evoli.pso_step
    @ (state: pso_state, inertia: f64, cognitive: f64, social: f64) -> pso_state
    + updates each particle's velocity and position using standard PSO dynamics
    ? v_new = inertia*v + cognitive*r1*(p_best - pos) + social*r2*(g_best - pos)
    # pso_step
    -> std.random.f64_unit
  evoli.pso_global_best
    @ (state: pso_state) -> tuple[list[f64], f64]
    + returns the swarm's best-ever position and its fitness
    # introspection
  evoli.ga_best
    @ (state: ga_state) -> tuple[list[f64], f64]
    + returns the fittest GA individual and its score
    # introspection
