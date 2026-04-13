# Requirement: "a library for converting between hexagonal grid cell indexes and GeoJSON geometries"

The hex grid and GeoJSON primitives are assumed available through std; the project layer handles the conversions and bulk set operations.

std
  std.geo
    std.geo.cell_to_boundary
      @ (cell: u64) -> list[coord]
      + returns the polygon boundary coordinates for a hex grid cell
      # geometry
    std.geo.coord_to_cell
      @ (lat: f64, lon: f64, resolution: i32) -> u64
      + returns the hex grid cell id containing the point at the given resolution
      # geometry
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document
      - returns error on invalid input
      # serialization
    std.json.encode_value
      @ (value: json_value) -> string
      + encodes a generic JSON value
      # serialization

hex_geojson
  hex_geojson.cell_to_feature
    @ (cell: u64, properties: map[string, string]) -> string
    + returns a GeoJSON Feature whose polygon is the cell boundary and whose properties include the cell id
    # conversion
    -> std.geo.cell_to_boundary
    -> std.json.encode_value
  hex_geojson.cells_to_feature_collection
    @ (cells: list[u64]) -> string
    + returns a GeoJSON FeatureCollection with one feature per cell
    # conversion
  hex_geojson.feature_to_cells
    @ (geojson: string, resolution: i32) -> result[list[u64], string]
    + returns the set of cells covering a GeoJSON polygon feature at the given resolution
    - returns error when the feature is not a polygon
    - returns error on invalid GeoJSON
    # conversion
    -> std.json.parse
    -> std.geo.coord_to_cell
  hex_geojson.polyfill_bbox
    @ (min_lat: f64, min_lon: f64, max_lat: f64, max_lon: f64, resolution: i32) -> list[u64]
    + returns all cells whose centers fall inside the bounding box
    # coverage
    -> std.geo.coord_to_cell
