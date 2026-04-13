# Requirement: "a library that detects image borders and converts them to GeoJSON"

Traces the non-transparent (or non-background) region of a raster image into a polygon, then georeferences and emits it as GeoJSON.

std
  std.image
    std.image.decode
      @ (data: bytes) -> result[bitmap_image, string]
      + decodes a PNG or JPEG buffer into a pixel bitmap
      - returns error on unsupported format
      # image
    std.image.pixel_at
      @ (img: bitmap_image, x: i32, y: i32) -> pixel_rgba
      + returns the RGBA value of a pixel
      # image
  std.json
    std.json.encode_value
      @ (value: json_value) -> string
      + serializes a JSON value to its textual form
      # serialization

borders
  borders.detect_mask
    @ (img: bitmap_image, background: pixel_rgba, tolerance: i32) -> boolean_mask
    + produces a boolean mask marking pixels that differ from the background within tolerance
    # segmentation
    -> std.image.pixel_at
  borders.trace_polygon
    @ (mask: boolean_mask) -> list[pair_i32]
    + walks the mask boundary and returns ordered pixel coordinates of the outer ring
    - returns empty list when the mask is fully empty
    # tracing
  borders.simplify_polygon
    @ (points: list[pair_i32], epsilon: f64) -> list[pair_i32]
    + reduces points using the Douglas-Peucker algorithm
    # simplification
  borders.pixel_to_geo
    @ (points: list[pair_i32], transform: affine_transform) -> list[pair_f64]
    + maps pixel coordinates to geographic coordinates using an affine transform
    # georeferencing
  borders.to_geojson
    @ (geo_points: list[pair_f64]) -> string
    + emits a GeoJSON Feature containing a single Polygon geometry
    # geojson
    -> std.json.encode_value
  borders.run
    @ (image_data: bytes, background: pixel_rgba, transform: affine_transform) -> result[string, string]
    + decodes the image, traces its border, and returns the GeoJSON representation
    - returns error when the image cannot be decoded
    # orchestration
    -> std.image.decode
