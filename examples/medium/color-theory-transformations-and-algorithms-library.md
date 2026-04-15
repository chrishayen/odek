# Requirement: "a color theory transformations and algorithms library"

Color-space conversions, differences, and blending. All pure functions.

std: (all units exist)

colour
  colour.rgb_to_hex
    fn (r: u8, g: u8, b: u8) -> string
    + returns "#RRGGBB" uppercase
    # conversion
  colour.hex_to_rgb
    fn (hex: string) -> result[tuple[u8, u8, u8], string]
    + parses "#RRGGBB" or "#RGB"
    - returns error on invalid length or non-hex characters
    # conversion
  colour.rgb_to_hsl
    fn (r: u8, g: u8, b: u8) -> tuple[f64, f64, f64]
    + returns (hue_degrees, saturation, lightness) with s and l in 0..1
    # conversion
  colour.hsl_to_rgb
    fn (h: f64, s: f64, l: f64) -> tuple[u8, u8, u8]
    + returns the RGB equivalent with channels clamped to 0..255
    # conversion
  colour.rgb_to_xyz
    fn (r: u8, g: u8, b: u8) -> tuple[f64, f64, f64]
    + returns CIE XYZ under the D65 illuminant
    # conversion
  colour.xyz_to_lab
    fn (x: f64, y: f64, z: f64) -> tuple[f64, f64, f64]
    + returns CIE L*a*b*
    # conversion
  colour.delta_e_cie76
    fn (lab_a: tuple[f64, f64, f64], lab_b: tuple[f64, f64, f64]) -> f64
    + returns the euclidean CIE76 color difference
    # comparison
  colour.delta_e_ciede2000
    fn (lab_a: tuple[f64, f64, f64], lab_b: tuple[f64, f64, f64]) -> f64
    + returns the CIEDE2000 color difference
    # comparison
  colour.mix
    fn (a: tuple[u8, u8, u8], b: tuple[u8, u8, u8], t: f64) -> tuple[u8, u8, u8]
    + returns the linear RGB interpolation between a and b at t in 0..1
    - returns a when t <= 0 and b when t >= 1
    # blending
  colour.luminance
    fn (r: u8, g: u8, b: u8) -> f64
    + returns the relative luminance per WCAG
    # measurement
  colour.contrast_ratio
    fn (a: tuple[u8, u8, u8], b: tuple[u8, u8, u8]) -> f64
    + returns the WCAG contrast ratio between two colors
    # measurement
    -> colour.luminance
