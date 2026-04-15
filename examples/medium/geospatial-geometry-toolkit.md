# Requirement: "a geospatial geometry toolkit"

Basic planar and great-circle operations on points, polygons, and lines.

std
  std.math
    std.math.sin
      fn (x: f64) -> f64
      + returns sine of x radians
      # math
    std.math.cos
      fn (x: f64) -> f64
      + returns cosine of x radians
      # math
    std.math.atan2
      fn (y: f64, x: f64) -> f64
      + returns the angle of the point (x, y) in radians
      # math
    std.math.sqrt
      fn (x: f64) -> f64
      + returns non-negative square root
      # math

geospatial
  geospatial.point_in_polygon
    fn (lat: f64, lon: f64, polygon: list[tuple[f64, f64]]) -> bool
    + returns true when the point lies strictly inside the polygon
    ? uses the ray-casting algorithm; treats the polygon as planar
    # containment
  geospatial.polygon_area_sq_meters
    fn (polygon: list[tuple[f64, f64]]) -> f64
    + returns the spherical area of the polygon on the earth ellipsoid in square meters
    ? uses the Green's theorem approximation for a WGS84 sphere
    # area
    -> std.math.sin
  geospatial.haversine_meters
    fn (lat_a: f64, lon_a: f64, lat_b: f64, lon_b: f64) -> f64
    + returns great-circle distance between two lat/lon points in meters
    # distance
    -> std.math.sin
    -> std.math.cos
    -> std.math.atan2
    -> std.math.sqrt
  geospatial.bearing_degrees
    fn (lat_a: f64, lon_a: f64, lat_b: f64, lon_b: f64) -> f64
    + returns the initial bearing from A to B in degrees clockwise from north
    # bearing
    -> std.math.sin
    -> std.math.cos
    -> std.math.atan2
  geospatial.bounding_box
    fn (points: list[tuple[f64, f64]]) -> result[tuple[f64, f64, f64, f64], string]
    + returns (min_lat, min_lon, max_lat, max_lon)
    - returns error when points is empty
    # bounding_box
  geospatial.simplify_line
    fn (points: list[tuple[f64, f64]], tolerance_meters: f64) -> list[tuple[f64, f64]]
    + returns a Douglas-Peucker-simplified polyline
    ? retains endpoints; removes points whose perpendicular distance to the simplified segment is below tolerance
    # simplification
  geospatial.centroid
    fn (polygon: list[tuple[f64, f64]]) -> result[tuple[f64, f64], string]
    + returns the area-weighted centroid of the polygon
    - returns error when polygon has fewer than three vertices
    # centroid
