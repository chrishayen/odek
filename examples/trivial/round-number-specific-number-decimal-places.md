# Requirement: "round a number to a specific number of decimal places"

Banker-free half-away-from-zero rounding at the requested precision.

std: (all units exist)

round_to
  round_to.round
    @ (value: f64, decimals: i32) -> f64
    + returns value rounded to the requested number of decimal places
    + rounds 1.234 with decimals=1 to 1.2
    + handles negative values symmetrically
    - treats negative decimals as rounding to the nearest power of ten
    # rounding
