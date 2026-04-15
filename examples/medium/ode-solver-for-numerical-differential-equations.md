# Requirement: "a library for numerically solving ordinary differential equations"

Integrators for first-order initial value problems, including a fixed-step Runge-Kutta 4 and an adaptive Dormand-Prince method.

std: (all units exist)

ode
  ode.new_problem
    fn (initial_state: list[f64], t0: f64, t_end: f64) -> problem
    + constructs an IVP description
    # construction
  ode.rk4_step
    fn (problem: problem, state: list[f64], t: f64, dt: f64, rhs_id: string) -> list[f64]
    + advances one Runge-Kutta 4 step and returns the new state
    # integration
  ode.rk4_solve
    fn (problem: problem, dt: f64, rhs_id: string) -> list[trajectory_point]
    + integrates from t0 to t_end with fixed step dt
    + returns the full trajectory
    # integration
  ode.dopri5_step
    fn (problem: problem, state: list[f64], t: f64, dt: f64, rhs_id: string) -> tuple[list[f64], f64]
    + advances one Dormand-Prince step and returns (new_state, error_estimate)
    # integration
  ode.dopri5_solve
    fn (problem: problem, tol: f64, rhs_id: string) -> list[trajectory_point]
    + integrates from t0 to t_end with adaptive step control
    + step size grows and shrinks to keep error below tol
    - aborts when step size underflows a minimum threshold
    # integration
  ode.interpolate
    fn (trajectory: list[trajectory_point], t: f64) -> optional[list[f64]]
    + returns state at t using linear interpolation between samples
    - returns none when t is outside the trajectory range
    # post_processing
