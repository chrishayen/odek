# Requirement: "a library of combinatorial operations: permutations, combinations, and combinations with replacement"

Generators that enumerate all arrangements of a small input list.

std: (all units exist)

combo
  combo.permutations
    @ (items: list[string], k: i32) -> list[list[string]]
    + returns all ordered selections of length k
    + returns a single empty arrangement when k is 0
    - returns empty when k exceeds the item count
    # enumeration
  combo.combinations
    @ (items: list[string], k: i32) -> list[list[string]]
    + returns all unordered selections of length k in input order
    - returns empty when k exceeds the item count
    # enumeration
  combo.combinations_with_replacement
    @ (items: list[string], k: i32) -> list[list[string]]
    + returns all unordered selections of length k allowing repeats
    # enumeration
  combo.count_permutations
    @ (n: i64, k: i64) -> i64
    + returns n! / (n-k)!
    - returns 0 when k > n
    # counting
  combo.count_combinations
    @ (n: i64, k: i64) -> i64
    + returns n choose k
    - returns 0 when k > n
    # counting
