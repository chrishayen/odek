# Requirement: "convert GeoJSON geometries into a covering of S2 cells"

Parses GeoJSON geometry, then asks an S2 indexer for a cell covering at a requested level range.

std
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses arbitrary JSON into a generic value tree
      - returns error on malformed input
      # serialization
  std.math
    std.math.deg_to_rad
      @ (degrees: f64) -> f64
      + converts degrees to radians
      # math

geos2
  geos2.parse_geometry
    @ (raw: string) -> result[geometry, string]
    + parses a Point, LineString, or Polygon GeoJSON feature
    - returns error when the "type" field is missing
    - returns error when coordinates are not numeric pairs
    # parsing
    -> std.json.parse_value
  geos2.latlng_to_cell
    @ (lat: f64, lng: f64, level: i32) -> u64
    + returns the cell id containing the given latitude/longitude at the requested level
    ? level is clamped to [0, 30]
    # projection
    -> std.math.deg_to_rad
  geos2.cover_geometry
    @ (geom: geometry, min_level: i32, max_level: i32, max_cells: i32) -> list[u64]
    + returns a point-cell list of length 1 for a Point
    + returns an ordered list of cells covering a LineString
    + returns an interior+boundary covering for a Polygon
    + honors the max_cells limit by merging to coarser parents
    # covering
    -> geos2.latlng_to_cell
  geos2.cell_to_latlng
    @ (cell: u64) -> tuple[f64, f64]
    + returns the center latitude/longitude of a cell
    # inspection
