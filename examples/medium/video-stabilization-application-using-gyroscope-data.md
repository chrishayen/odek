# Requirement: "a video stabilization library using gyroscope data"

Fuses per-frame timestamps with gyroscope samples to compute stabilizing rotations and warp parameters per frame.

std
  std.math
    std.math.interpolate_linear
      @ (x: f64, x0: f64, x1: f64, y0: f64, y1: f64) -> f64
      + returns linearly interpolated y for x in [x0, x1]
      # math
    std.math.quaternion_multiply
      @ (a: quaternion, b: quaternion) -> quaternion
      + returns the hamilton product a*b
      # math
    std.math.quaternion_normalize
      @ (q: quaternion) -> quaternion
      + returns q scaled to unit length
      - returns identity when magnitude is zero
      # math

stabilizer
  stabilizer.parse_gyro_samples
    @ (raw: bytes) -> result[list[gyro_sample], string]
    + parses timestamped angular velocity samples (x, y, z)
    - returns error on malformed sample record
    # input
  stabilizer.integrate_to_orientations
    @ (samples: list[gyro_sample]) -> list[orientation]
    + integrates angular velocities over time into per-sample quaternion orientations
    + returns [] when samples is empty
    # motion
    -> std.math.quaternion_multiply
    -> std.math.quaternion_normalize
  stabilizer.orientation_at
    @ (orientations: list[orientation], timestamp_ns: i64) -> optional[quaternion]
    + returns the interpolated orientation at the given timestamp
    - returns none when timestamp is outside the sample range
    # motion
    -> std.math.interpolate_linear
  stabilizer.smooth_orientations
    @ (orientations: list[orientation], window_ns: i64) -> list[orientation]
    + applies a time-windowed average producing a smoothed motion path
    # motion
  stabilizer.compute_correction
    @ (raw: quaternion, smoothed: quaternion) -> quaternion
    + returns the quaternion that rotates raw onto smoothed
    # correction
    -> std.math.quaternion_multiply
  stabilizer.frame_transform
    @ (correction: quaternion, focal_length_px: f64) -> transform_matrix
    + converts a 3D correction into a 2D warp matrix for the given focal length
    # projection
