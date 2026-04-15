# Requirement: "a collection of classic robotics algorithms: path planning, localization, and control"

Covers three subsystems a robotics library typically bundles. Math primitives used by all of them live in std.

std
  std.math
    std.math.hypot
      fn (dx: f64, dy: f64) -> f64
      + returns sqrt(dx*dx + dy*dy) without overflow for large magnitudes
      # math
    std.math.wrap_angle
      fn (radians: f64) -> f64
      + wraps an angle into the range (-pi, pi]
      # math
  std.random
    std.random.gaussian
      fn (mean: f64, stddev: f64) -> f64
      + returns a sample from a normal distribution
      # random

robotics
  robotics.grid_new
    fn (width: i32, height: i32) -> grid_map
    + creates an empty occupancy grid with all cells free
    # map
  robotics.grid_set_obstacle
    fn (grid: grid_map, x: i32, y: i32) -> grid_map
    + marks a cell as blocked
    # map
  robotics.astar
    fn (grid: grid_map, start_xy: tuple[i32, i32], goal_xy: tuple[i32, i32]) -> result[list[tuple[i32, i32]], string]
    + returns a path from start to goal using 8-connected A* with Euclidean heuristic
    - returns error when no path exists
    # path_planning
    -> std.math.hypot
  robotics.rrt
    fn (grid: grid_map, start_xy: tuple[f64, f64], goal_xy: tuple[f64, f64], max_iters: i32) -> result[list[tuple[f64, f64]], string]
    + returns a path from a randomly grown rapidly-exploring tree
    - returns error when the iteration budget is exhausted without reaching the goal
    # path_planning
    -> std.math.hypot
    -> std.random.gaussian
  robotics.pure_pursuit_step
    fn (path: list[tuple[f64, f64]], pose: pose2d, lookahead: f64) -> f64
    + returns the desired steering angle for the given pose and lookahead distance
    # control
    -> std.math.wrap_angle
  robotics.pid_new
    fn (kp: f64, ki: f64, kd: f64) -> pid_state
    + creates a PID controller with the given gains and zero history
    # control
  robotics.pid_step
    fn (state: pid_state, error: f64, dt: f64) -> tuple[pid_state, f64]
    + advances the controller and returns the next control output
    # control
  robotics.ekf_predict
    fn (state: ekf_state, control: list[f64], dt: f64) -> ekf_state
    + advances an extended Kalman filter by one motion step
    # localization
  robotics.ekf_update
    fn (state: ekf_state, observation: list[f64]) -> ekf_state
    + applies a measurement update to the filter state
    # localization
  robotics.particle_filter_step
    fn (particles: list[particle], control: list[f64], observation: list[f64]) -> list[particle]
    + predicts, weights, and resamples a particle set in one step
    # localization
    -> std.random.gaussian
  robotics.icp_align
    fn (source: list[tuple[f64, f64]], target: list[tuple[f64, f64]], max_iters: i32) -> tuple[f64, f64, f64]
    + returns (dx, dy, dtheta) that best aligns source to target
    - returns zeros when either cloud is empty
    # scan_matching
    -> std.math.hypot
    -> std.math.wrap_angle
