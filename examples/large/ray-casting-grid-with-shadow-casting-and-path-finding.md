# Requirement: "a 2D grid with ray casting, shadow casting, and path finding"

A tile grid supporting line-of-sight queries, field-of-view computation, and shortest-path search.

std
  std.collections
    std.collections.new_priority_queue
      @ () -> priority_queue
      + creates an empty min-priority queue keyed by i32 cost
      # collections
    std.collections.pq_push
      @ (q: priority_queue, key: i32, value: i32) -> priority_queue
      + inserts value with the given key
      # collections
    std.collections.pq_pop_min
      @ (q: priority_queue) -> optional[tuple[i32, i32, priority_queue]]
      + removes and returns the minimum-key entry as (key, value, remaining)
      + returns none when empty
      # collections

grid
  grid.new
    @ (width: i32, height: i32) -> grid_state
    + creates a grid of the given size with every cell passable
    # construction
  grid.set_blocked
    @ (g: grid_state, x: i32, y: i32, blocked: bool) -> grid_state
    + marks the cell at (x, y) as blocking or passable
    # mutation
  grid.is_blocked
    @ (g: grid_state, x: i32, y: i32) -> bool
    + returns true when the cell is out of bounds or marked blocked
    # query
  grid.cast_ray
    @ (g: grid_state, x0: i32, y0: i32, x1: i32, y1: i32) -> list[tuple[i32, i32]]
    + returns the cells traversed from the source to the first blocked cell or the endpoint, using a Bresenham-style step
    # ray_casting
    -> grid.is_blocked
  grid.line_of_sight
    @ (g: grid_state, x0: i32, y0: i32, x1: i32, y1: i32) -> bool
    + returns true when no blocked cell lies between the endpoints
    # ray_casting
    -> grid.cast_ray
  grid.field_of_view
    @ (g: grid_state, cx: i32, cy: i32, radius: i32) -> list[tuple[i32, i32]]
    + returns every visible cell within radius of the source using recursive shadow casting
    # shadow_casting
    -> grid.is_blocked
  grid.neighbors
    @ (g: grid_state, x: i32, y: i32) -> list[tuple[i32, i32]]
    + returns the up to four orthogonal passable neighbors of (x, y)
    # query
    -> grid.is_blocked
  grid.manhattan
    @ (ax: i32, ay: i32, bx: i32, by: i32) -> i32
    + returns the manhattan distance between two cells
    # heuristic
  grid.find_path
    @ (g: grid_state, sx: i32, sy: i32, gx: i32, gy: i32) -> optional[list[tuple[i32, i32]]]
    + returns the shortest cell path from source to goal using A* with a manhattan heuristic
    - returns none when no path exists
    # path_finding
    -> grid.neighbors
    -> grid.manhattan
    -> std.collections.new_priority_queue
    -> std.collections.pq_push
    -> std.collections.pq_pop_min
