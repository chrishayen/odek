# Requirement: "a symbolic mathematics library"

Expressions are an opaque tree; the project exposes construction, simplification, and differentiation.

std: (all units exist)

symmath
  symmath.constant
    fn (value: f64) -> expr
    + returns an expression node representing a numeric constant
    # construction
  symmath.variable
    fn (name: string) -> expr
    + returns an expression node representing a named variable
    # construction
  symmath.add
    fn (left: expr, right: expr) -> expr
    + returns an expression representing left + right
    # construction
  symmath.multiply
    fn (left: expr, right: expr) -> expr
    + returns an expression representing left * right
    # construction
  symmath.power
    fn (base: expr, exponent: expr) -> expr
    + returns an expression representing base raised to exponent
    # construction
  symmath.simplify
    fn (e: expr) -> expr
    + folds constant arithmetic subtrees
    + removes zero terms in sums and unit factors in products
    + collapses x^0 to 1 and x^1 to x
    # simplification
  symmath.differentiate
    fn (e: expr, variable: string) -> expr
    + returns the symbolic derivative with respect to the given variable
    + applies sum, product, chain, and power rules
    # differentiation
  symmath.evaluate
    fn (e: expr, bindings: map[string, f64]) -> result[f64, string]
    + returns the numeric value under the given variable bindings
    - returns error when a free variable is unbound
    # evaluation
  symmath.to_string
    fn (e: expr) -> string
    + returns a parenthesized textual form of the expression
    # serialization
