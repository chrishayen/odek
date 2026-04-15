# Requirement: "an astronomy computation library"

A focused subset: coordinate conversions, unit handling for angles, and an apparent-magnitude helper. Generic math lives in std.

std
  std.math
    std.math.deg_to_rad
      fn (deg: f64) -> f64
      + converts degrees to radians
      # math
    std.math.rad_to_deg
      fn (rad: f64) -> f64
      + converts radians to degrees
      # math
    std.math.atan2
      fn (y: f64, x: f64) -> f64
      + returns the angle in radians between the positive x-axis and (x, y)
      # math

astronomy
  astronomy.parse_hms
    fn (hms: string) -> result[f64, string]
    + parses a right ascension string "HH:MM:SS.s" into degrees
    - returns error when the format is malformed
    # parsing
  astronomy.parse_dms
    fn (dms: string) -> result[f64, string]
    + parses a declination string "+DD:MM:SS.s" into degrees
    - returns error when the format is malformed
    # parsing
  astronomy.equatorial_to_horizontal
    fn (ra_deg: f64, dec_deg: f64, lat_deg: f64, lst_deg: f64) -> tuple[f64, f64]
    + returns (azimuth_deg, altitude_deg) for the observer latitude and local sidereal time
    # coordinates
    -> std.math.deg_to_rad
    -> std.math.rad_to_deg
    -> std.math.atan2
  astronomy.angular_separation
    fn (ra1_deg: f64, dec1_deg: f64, ra2_deg: f64, dec2_deg: f64) -> f64
    + returns the great-circle separation in degrees between two celestial positions
    # coordinates
    -> std.math.deg_to_rad
    -> std.math.rad_to_deg
  astronomy.apparent_magnitude
    fn (absolute_mag: f64, distance_parsecs: f64) -> result[f64, string]
    + returns apparent magnitude given absolute magnitude and distance in parsecs
    - returns error when distance_parsecs is not positive
    # photometry
