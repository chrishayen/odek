# Requirement: "a function to add two integers"

One arithmetic operation. No helpers, no variants, no std — integer addition is a language primitive.

std: (all units exist)

arithmetic
  arithmetic.add
    @ (a: i32, b: i32) -> i32
    + returns 5 when given 2 and 3
    + returns 0 when given 0 and 0
    + returns -5 when given -2 and -3
    - overflow behavior is undefined; caller must ensure the result fits in i32
    # arithmetic
