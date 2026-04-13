# Requirement: "a spatial data and geometric algorithms library"

Core 2D geometry types, predicates, and operations over points, lines, and polygons.

std
  std.math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the non-negative square root
      # math
    std.math.abs
      @ (x: f64) -> f64
      + returns the absolute value
      # math

geometry
  geometry.point
    @ (x: f64, y: f64) -> point
    + builds a 2D point with the given coordinates
    # construction
  geometry.distance
    @ (a: point, b: point) -> f64
    + returns the euclidean distance between two points
    + returns 0 when both points are identical
    # measurement
    -> std.math.sqrt
  geometry.segment_intersects
    @ (a1: point, a2: point, b1: point, b2: point) -> bool
    + returns true when two segments cross or touch
    - returns false for disjoint segments
    ? collinear overlapping segments count as intersecting
    # predicate
  geometry.polygon_area
    @ (vertices: list[point]) -> f64
    + returns the unsigned area using the shoelace formula
    + returns 0 for polygons with fewer than three vertices
    # measurement
    -> std.math.abs
  geometry.polygon_contains_point
    @ (vertices: list[point], p: point) -> bool
    + returns true when the point lies inside the polygon
    - returns false for points strictly outside
    ? uses the ray-casting algorithm; boundary inclusion is implementation-defined
    # predicate
  geometry.bounding_box
    @ (points: list[point]) -> result[bounding_box, string]
    + returns the axis-aligned minimum bounding box
    - returns error when the list is empty
    # measurement
  geometry.convex_hull
    @ (points: list[point]) -> list[point]
    + returns the convex hull vertices in counter-clockwise order
    + returns the input unchanged when three or fewer points are given
    # hull
  geometry.polygon_union
    @ (a: list[point], b: list[point]) -> result[list[list[point]], string]
    + returns the polygons forming the union of two simple polygons
    - returns error when either input is not a simple polygon
    # boolean_op
  geometry.polygon_intersection
    @ (a: list[point], b: list[point]) -> result[list[list[point]], string]
    + returns the polygons forming the intersection
    + returns an empty list when the inputs are disjoint
    - returns error when either input is not a simple polygon
    # boolean_op
  geometry.buffer
    @ (vertices: list[point], distance: f64) -> list[point]
    + returns a polygon expanded outward by the given distance
    + negative distance shrinks the polygon inward
    # offsetting
