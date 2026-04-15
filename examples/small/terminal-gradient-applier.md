# Requirement: "applies smooth color gradients to terminal text"

Interpolates between stop colors and wraps each character in an ANSI truecolor escape.

std: (all units exist)

gradient
  gradient.parse_hex
    fn (hex: string) -> result[rgb_color, string]
    + parses "#rrggbb" or "rrggbb" into an rgb triple
    - returns error on wrong length or non-hex characters
    # parsing
  gradient.interpolate
    fn (a: rgb_color, b: rgb_color, t: f64) -> rgb_color
    + returns the linear interpolation between a and b at t in [0, 1]
    ? t outside [0, 1] is clamped
    # math
  gradient.apply
    fn (text: string, stops: list[rgb_color]) -> string
    + returns text with each visible character wrapped in an ANSI 24-bit color escape
    + resets color at the end with the standard reset sequence
    - returns text unchanged when stops is empty
    ? whitespace characters receive colors too but are visually invisible
    # rendering
