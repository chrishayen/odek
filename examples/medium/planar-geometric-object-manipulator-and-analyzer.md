# Requirement: "a library for manipulation and analysis of planar geometric objects"

Points, line segments, and polygons in the Cartesian plane with the standard predicates and measurements.

std: (all units exist)

geometry
  geometry.point
    fn (x: f64, y: f64) -> point
    + constructs a point from coordinates
    # construction
  geometry.distance
    fn (a: point, b: point) -> f64
    + returns the Euclidean distance between two points
    # measurement
  geometry.segment
    fn (a: point, b: point) -> segment
    + constructs a segment from two endpoints
    # construction
  geometry.segments_intersect
    fn (s1: segment, s2: segment) -> bool
    + returns true when segments share at least one point
    - returns false for strictly disjoint segments
    # predicates
  geometry.polygon
    fn (vertices: list[point]) -> result[polygon, string]
    + constructs a simple polygon from at least three vertices
    - returns error when vertices has fewer than three points
    # construction
  geometry.polygon_area
    fn (p: polygon) -> f64
    + returns the unsigned area via the shoelace formula
    # measurement
  geometry.point_in_polygon
    fn (pt: point, p: polygon) -> bool
    + returns true when the point lies inside the polygon using ray casting
    - returns false for points strictly outside
    ? points on the boundary are reported as inside
    # predicates
  geometry.bounding_box
    fn (p: polygon) -> tuple[point, point]
    + returns the lower-left and upper-right corners of the axis-aligned bounding box
    # measurement
  geometry.convex_hull
    fn (points: list[point]) -> result[polygon, string]
    + returns the convex hull of a set of points using the monotone chain algorithm
    - returns error when fewer than three unique points are provided
    # analysis
