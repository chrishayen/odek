# Requirement: "a library for tiled spatial image survey projection and lookup"

Given a celestial coordinate and zoom level, computes which tile identifiers cover the viewport. Supports multiple projections.

std: (all units exist)

sky_tiles
  sky_tiles.lonlat_to_unit_vector
    fn (lon_deg: f64, lat_deg: f64) -> list[f64]
    + returns a unit vector on the sphere for the given longitude/latitude
    # math
  sky_tiles.healpix_index
    fn (lon_deg: f64, lat_deg: f64, order: i32) -> i64
    + returns the HEALPix pixel index at the given order
    - returns -1 for invalid coordinates
    # indexing
  sky_tiles.tiles_in_view
    fn (center_lon: f64, center_lat: f64, radius_deg: f64, order: i32) -> list[i64]
    + returns the list of tile indices overlapping the circular viewport
    # query
  sky_tiles.project_gnomonic
    fn (lon_deg: f64, lat_deg: f64, center_lon: f64, center_lat: f64) -> tuple[f64, f64]
    + returns gnomonic (x, y) projection coordinates
    # projection
  sky_tiles.project_mercator
    fn (lon_deg: f64, lat_deg: f64) -> tuple[f64, f64]
    + returns mercator (x, y) coordinates
    - returns large values approaching the poles
    # projection
  sky_tiles.tile_url_path
    fn (tile_index: i64, order: i32) -> string
    + returns the conventional tile path segment for the given index
    # formatting
