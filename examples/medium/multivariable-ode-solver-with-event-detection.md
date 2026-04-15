# Requirement: "a multivariable ODE solver with event detection"

Integrates a user-supplied derivative function using fixed-step RK4 and detects events where a caller-supplied scalar function crosses zero.

std: (all units exist)

odesolve
  odesolve.new_system
    fn (initial_state: list[f64], t0: f64) -> ode_state
    + creates a system with the given initial variable values at time t0
    # construction
  odesolve.step_rk4
    fn (state: ode_state, dt: f64, deriv: list[f64]) -> ode_state
    + advances the state by dt using classical RK4
    ? the caller supplies the derivative vector; the library does not own the physics
    # integration
  odesolve.integrate
    fn (state: ode_state, dt: f64, steps: i32, deriv_seq: list[list[f64]]) -> ode_state
    + advances the state for the given number of steps using pre-computed derivatives
    - returns the starting state unchanged when steps is zero
    # integration
  odesolve.detect_zero_crossing
    fn (prev_value: f64, curr_value: f64, prev_t: f64, curr_t: f64) -> optional[f64]
    + returns the linearly interpolated time at which the value crosses zero
    - returns none when both samples have the same sign
    # events
  odesolve.record_event
    fn (state: ode_state, t: f64, tag: string) -> ode_state
    + appends an event with its time and tag to the state
    # events
  odesolve.trajectory
    fn (state: ode_state) -> list[tuple[f64, list[f64]]]
    + returns the time and variable history recorded so far
    # query
