# Requirement: "a root-finding library for quadratic functions"

Solves ax^2 + bx + c = 0 over the reals.

std: (all units exist)

quadratic
  quadratic.roots
    fn (a: f64, b: f64, c: f64) -> result[list[f64], string]
    + returns both roots when the discriminant is positive
    + returns a single root when the discriminant is zero
    - returns an empty list when the discriminant is negative
    - returns error when a is zero (not a quadratic)
    # root_finding
