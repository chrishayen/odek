# Requirement: "a library that generates SVG patterns from an input string"

Hash the input to seed deterministic geometry and color choices, then emit tiled SVG shapes.

std
  std.hash
    std.hash.sha1
      @ (data: bytes) -> bytes
      + returns the SHA-1 digest of data (20 bytes)
      # hashing
  std.xml
    std.xml.escape_attr
      @ (value: string) -> string
      + escapes &, <, >, and quotes for XML attribute values
      # serialization

geopattern
  geopattern.seed_from_string
    @ (input: string) -> bytes
    + returns a 20-byte hash to drive pattern selection and colors
    # seeding
    -> std.hash.sha1
  geopattern.pick_style
    @ (seed: bytes) -> pattern_style
    + chooses one of the built-in styles (hexagons, triangles, squares, chevrons, xes, bricks, plus_signs)
    # selection
  geopattern.pick_colors
    @ (seed: bytes) -> tuple[string, string]
    + returns (background, foreground) hex colors derived from the seed
    # color
  geopattern.render_tile
    @ (style: pattern_style, seed: bytes, fg: string) -> string
    + emits the SVG group for a single tile of the chosen style
    # rendering
    -> std.xml.escape_attr
  geopattern.generate
    @ (input: string) -> string
    + returns a complete SVG document tiling the generated pattern
    # pipeline
