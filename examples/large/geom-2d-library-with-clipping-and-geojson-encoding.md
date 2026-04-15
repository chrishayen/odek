# Requirement: "a 2D geometry library with clipping and GeoJSON encoding"

Core geometric types plus a polygon clipper and a GeoJSON round-trip. Tile encoding is out of scope — it is a separate concern.

std
  std.json
    std.json.parse_value
      fn (raw: string) -> result[json_value, string]
      + parses any JSON value into a tagged tree
      - returns error on malformed JSON
      # serialization
    std.json.encode_value
      fn (v: json_value) -> string
      + encodes a tagged JSON tree to a string
      # serialization
  std.math
    std.math.min_f64
      fn (a: f64, b: f64) -> f64
      + returns the smaller of two values
      # math
    std.math.max_f64
      fn (a: f64, b: f64) -> f64
      + returns the larger of two values
      # math

geom
  geom.point_new
    fn (x: f64, y: f64) -> point
    + constructs a 2D point
    # construction
  geom.bbox_of_points
    fn (pts: list[point]) -> optional[bbox]
    + returns the axis-aligned bounding box containing all points
    - returns none for an empty list
    # bbox
    -> std.math.min_f64
    -> std.math.max_f64
  geom.segment_intersect
    fn (a0: point, a1: point, b0: point, b1: point) -> optional[point]
    + returns the intersection point when two segments cross
    - returns none when segments are parallel or disjoint
    # intersection
  geom.polygon_area
    fn (ring: list[point]) -> f64
    + returns signed area via the shoelace formula
    ? positive for counter-clockwise rings
    # area
  geom.polygon_contains
    fn (ring: list[point], p: point) -> bool
    + returns true when the point lies strictly inside the ring
    + uses the ray-casting rule
    - returns false for points outside or exactly on the boundary
    # point_in_polygon
  geom.clip_polygon_to_bbox
    fn (ring: list[point], box: bbox) -> list[point]
    + returns the Sutherland-Hodgman clip of the ring against the bbox
    + returns an empty ring when the polygon is entirely outside the box
    ? input ring is assumed closed (first and last point equal)
    # clipping
    -> geom.segment_intersect
  geom.line_clip_to_bbox
    fn (line: list[point], box: bbox) -> list[list[point]]
    + returns the subsegments of the line that lie inside the bbox
    + splits into multiple pieces when the line re-enters the box
    # clipping
  geom.from_geojson
    fn (raw: string) -> result[geometry, string]
    + parses a GeoJSON Point, LineString, Polygon, or their Multi variants
    - returns error for unsupported types or malformed coordinates
    # geojson_decode
    -> std.json.parse_value
  geom.to_geojson
    fn (g: geometry) -> string
    + encodes a geometry as a minimal GeoJSON object
    # geojson_encode
    -> std.json.encode_value
