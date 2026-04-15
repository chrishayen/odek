# Requirement: "a library that generates a harmonious 16-color terminal palette from a single seed (color or string)"

Derives a base hue from a seed, then constructs 16 terminal slots by rotating and varying lightness.

std
  std.hash
    std.hash.fnv32
      fn (data: string) -> u32
      + deterministic 32-bit hash of the input
      # hashing

chroma16
  chroma16.seed_to_hue
    fn (seed: string) -> f64
    + returns a hue in [0, 360) derived from the seed; accepts "#rrggbb" literally or hashes arbitrary strings
    # seeding
    -> std.hash.fnv32
  chroma16.hsl_to_hex
    fn (h: f64, s: f64, l: f64) -> string
    + converts HSL to a "#rrggbb" string
    + clamps saturation and lightness to [0, 1]
    # color_space
  chroma16.build_palette
    fn (seed: string) -> list[string]
    + returns exactly 16 hex colors: background, foreground, 8 normal ansi colors, 8 bright variants
    + is deterministic for a given seed
    # palette
