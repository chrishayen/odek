# Requirement: "convert between celsius and fahrenheit"

Two pure functions. The bodies are one-line formulas — nothing to factor out, nothing reusable enough for std.

std: (all units exist)

temperature
  temperature.celsius_to_fahrenheit
    fn (c: f64) -> f64
    + returns 32.0 when given 0.0
    + returns 212.0 when given 100.0
    + returns -40.0 when given -40.0 (the only fixed point)
    # conversion
  temperature.fahrenheit_to_celsius
    fn (f: f64) -> f64
    + returns 0.0 when given 32.0
    + returns 100.0 when given 212.0
    # conversion
