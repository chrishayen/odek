# Requirement: "a modular geospatial processing and analysis library"

Core geometry and analysis primitives over longitude/latitude features. All distances are great-circle on a sphere; planar operations are explicit.

std
  std.math
    std.math.sin
      fn (x: f64) -> f64
      + returns the sine of x in radians
      # math
    std.math.cos
      fn (x: f64) -> f64
      + returns the cosine of x in radians
      # math
    std.math.atan2
      fn (y: f64, x: f64) -> f64
      + returns the angle of (x, y) in radians
      # math
    std.math.sqrt
      fn (x: f64) -> f64
      + returns the non-negative square root
      - returns NaN when x is negative
      # math
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document
      - returns error on invalid syntax
      # serialization
    std.json.encode
      fn (v: json_value) -> string
      + serializes a json value
      # serialization

geo
  geo.point
    fn (lon: f64, lat: f64) -> point
    + constructs a geographic point
    - rejects coordinates outside [-180, 180] x [-90, 90]
    # construction
  geo.distance
    fn (a: point, b: point) -> f64
    + returns great-circle distance in meters using the haversine formula
    # distance
    -> std.math.sin
    -> std.math.cos
    -> std.math.atan2
    -> std.math.sqrt
  geo.bearing
    fn (a: point, b: point) -> f64
    + returns initial bearing from a to b in degrees [0, 360)
    # bearing
    -> std.math.sin
    -> std.math.cos
    -> std.math.atan2
  geo.bbox_of
    fn (points: list[point]) -> bbox
    + returns the minimum bounding box containing all points
    - returns a zero bbox when the list is empty
    # bounding
  geo.bbox_contains
    fn (box: bbox, p: point) -> bool
    + returns true when the point lies inside the box inclusive of edges
    # spatial_query
  geo.polygon
    fn (ring: list[point]) -> result[polygon, string]
    + constructs a simple polygon from a closed ring
    - returns error when the ring has fewer than four vertices or is not closed
    # construction
  geo.polygon_area
    fn (poly: polygon) -> f64
    + returns the spherical excess area in square meters
    # area
    -> std.math.sin
  geo.polygon_contains
    fn (poly: polygon, p: point) -> bool
    + returns true when the point is inside via winding number
    # spatial_query
  geo.line_length
    fn (line: list[point]) -> f64
    + returns total great-circle length of the polyline in meters
    # distance
    -> geo.distance
  geo.simplify
    fn (line: list[point], tolerance_m: f64) -> list[point]
    + returns a reduced polyline using Douglas-Peucker with the given tolerance
    # simplification
    -> geo.distance
  geo.feature_to_geojson
    fn (p: polygon) -> string
    + returns a GeoJSON Feature wrapping the polygon
    # serialization
    -> std.json.encode
  geo.feature_from_geojson
    fn (raw: string) -> result[polygon, string]
    + parses a GeoJSON Polygon Feature
    - returns error when geometry type is not Polygon
    # serialization
    -> std.json.parse
