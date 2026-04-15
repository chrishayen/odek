# Requirement: "a geographic data frame library"

A tabular structure where one column is a geometry. Geometry primitives live in std.

std
  std.geometry
    std.geometry.point
      fn (x: f64, y: f64) -> geometry
      + constructs a point geometry
      # geometry
    std.geometry.bbox
      fn (g: geometry) -> tuple[f64, f64, f64, f64]
      + returns (min_x, min_y, max_x, max_y)
      # geometry
    std.geometry.contains
      fn (outer: geometry, inner: geometry) -> bool
      + returns true when outer fully contains inner
      # geometry
    std.geometry.distance
      fn (a: geometry, b: geometry) -> f64
      + returns the planar distance between two geometries
      # geometry

geoframe
  geoframe.new
    fn (columns: list[string]) -> geoframe_state
    + creates an empty frame with the given non-geometry column names
    ? the geometry column is implicit and always named "geometry"
    # construction
  geoframe.add_row
    fn (frame: geoframe_state, values: list[string], geom: geometry) -> geoframe_state
    + appends a row aligning values to the declared columns
    - returns unchanged when values count does not match columns
    # mutation
  geoframe.filter_within
    fn (frame: geoframe_state, region: geometry) -> geoframe_state
    + returns a new frame containing rows whose geometry is fully inside region
    # spatial_query
    -> std.geometry.contains
  geoframe.nearest
    fn (frame: geoframe_state, target: geometry, k: i32) -> list[i32]
    + returns row indices of the k geometries closest to target
    # spatial_query
    -> std.geometry.distance
  geoframe.total_bounds
    fn (frame: geoframe_state) -> tuple[f64, f64, f64, f64]
    + returns the bounding box covering every row's geometry
    # aggregation
    -> std.geometry.bbox
  geoframe.select
    fn (frame: geoframe_state, column: string) -> list[string]
    + returns the values of one non-geometry column in row order
    # access
  geoframe.length
    fn (frame: geoframe_state) -> i32
    + returns the number of rows
    # access
