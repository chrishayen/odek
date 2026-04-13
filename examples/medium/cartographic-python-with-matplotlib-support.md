# Requirement: "a cartographic library for map projection and coordinate transformation"

Converts between geographic coordinates and projected plane coordinates, and produces projected polyline data for plotting.

std: (all units exist)

cartography
  cartography.wgs84_to_mercator
    @ (lat: f64, lon: f64) -> tuple[f64, f64]
    + projects a WGS84 lat/lon to Web Mercator x/y in meters
    + returns (0, 0) for (0, 0)
    # projection
  cartography.mercator_to_wgs84
    @ (x: f64, y: f64) -> tuple[f64, f64]
    + inverse of wgs84_to_mercator
    # projection
  cartography.project_equirect
    @ (lat: f64, lon: f64, center_lat: f64, center_lon: f64) -> tuple[f64, f64]
    + projects with the equirectangular projection about the given center
    # projection
  cartography.haversine_km
    @ (lat1: f64, lon1: f64, lat2: f64, lon2: f64) -> f64
    + returns the great-circle distance between two points in kilometers
    # geometry
  cartography.project_polyline
    @ (points: list[tuple[f64,f64]], projection: projection) -> list[tuple[f64,f64]]
    + applies a named projection to each lat/lon point in order
    # polyline
  cartography.densify_polyline
    @ (points: list[tuple[f64,f64]], max_segment_km: f64) -> list[tuple[f64,f64]]
    + subdivides long edges so no segment exceeds max_segment_km
    ? prevents straight lines on the plane from looking wrong on curved projections
    # polyline
  cartography.bounds
    @ (points: list[tuple[f64,f64]]) -> tuple[f64,f64,f64,f64]
    + returns (min_x, min_y, max_x, max_y) for the projected points
    # geometry
