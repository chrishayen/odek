# Requirement: "a geodetic coordinate conversion and transformation library"

Convert between geographic (lat/lon), geocentric (ECEF), and projected (UTM, Web Mercator) coordinates using named ellipsoids and Helmert transforms.

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
    std.math.tan
      fn (x: f64) -> f64
      + returns the tangent of x in radians
      # math
    std.math.sqrt
      fn (x: f64) -> f64
      + returns the square root of x
      - returns NaN for negative x
      # math
    std.math.atan2
      fn (y: f64, x: f64) -> f64
      + returns the angle whose tangent is y/x, in the correct quadrant
      # math

geodesy
  geodesy.ellipsoid
    fn (name: string) -> result[ellipsoid, string]
    + returns the ellipsoid parameters for wgs84, grs80, airy1830, clarke1866
    - returns error for unknown ellipsoid names
    # reference_data
  geodesy.geographic_to_ecef
    fn (e: ellipsoid, lat: f64, lon: f64, height: f64) -> point3
    + converts geographic lat/lon/height to earth-centered earth-fixed XYZ
    # conversion
    -> std.math.sin
    -> std.math.cos
    -> std.math.sqrt
  geodesy.ecef_to_geographic
    fn (e: ellipsoid, p: point3) -> tuple[f64, f64, f64]
    + returns (lat, lon, height) iterating to centimeter precision
    # conversion
    -> std.math.atan2
    -> std.math.sqrt
  geodesy.helmert_transform
    fn (p: point3, params: helmert_params) -> point3
    + applies a 7-parameter Helmert datum transform
    # datum_shift
  geodesy.geographic_to_utm
    fn (lat: f64, lon: f64) -> tuple[i32, f64, f64]
    + returns (zone, easting, northing) using the WGS84 ellipsoid
    - returns zone 0 when latitude is outside the UTM valid range
    # projection
    -> std.math.sin
    -> std.math.cos
    -> std.math.tan
  geodesy.utm_to_geographic
    fn (zone: i32, easting: f64, northing: f64, northern: bool) -> tuple[f64, f64]
    + returns (lat, lon) in degrees
    # projection
    -> std.math.sin
    -> std.math.cos
  geodesy.geographic_to_web_mercator
    fn (lat: f64, lon: f64) -> tuple[f64, f64]
    + returns (x, y) in EPSG:3857 meters
    - clamps latitudes beyond +-85.05 to avoid infinity
    # projection
    -> std.math.tan
  geodesy.web_mercator_to_geographic
    fn (x: f64, y: f64) -> tuple[f64, f64]
    + returns (lat, lon) in degrees from EPSG:3857 meters
    # projection
    -> std.math.atan2
