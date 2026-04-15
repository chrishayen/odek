# Requirement: "a quadtree spatial index"

A point quadtree over a bounded 2D region. Supports insertion and range queries.

std: (all units exist)

quadtree
  quadtree.new
    fn (min_x: f64, min_y: f64, max_x: f64, max_y: f64, capacity: i32) -> quadtree_state
    + creates an empty quadtree covering the given axis-aligned bounds
    ? capacity is the number of points a node holds before subdividing
    # construction
  quadtree.insert
    fn (state: quadtree_state, id: string, x: f64, y: f64) -> result[quadtree_state, string]
    + inserts a point with the given id and returns the updated tree
    - returns error when the point is outside the root bounds
    # insertion
  quadtree.subdivide
    fn (state: quadtree_state, node: i32) -> quadtree_state
    + splits a leaf node into four quadrant children, redistributing its points
    ? called internally when a leaf exceeds capacity
    # subdivide
  quadtree.query_range
    fn (state: quadtree_state, min_x: f64, min_y: f64, max_x: f64, max_y: f64) -> list[string]
    + returns ids of all points inside the given axis-aligned rectangle
    + prunes subtrees whose bounds do not intersect the query
    # range_query
  quadtree.query_radius
    fn (state: quadtree_state, cx: f64, cy: f64, radius: f64) -> list[string]
    + returns ids of all points within Euclidean distance `radius` of (cx, cy)
    # radius_query
  quadtree.size
    fn (state: quadtree_state) -> i32
    + returns the total number of points stored
    # introspection
