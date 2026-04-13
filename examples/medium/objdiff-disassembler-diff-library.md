# Requirement: "a diffing library for comparing two versions of a disassembled program"

Compares named functions between two object snapshots and reports per-function similarity.

std: (all units exist)

objdiff
  objdiff.new
    @ () -> diff_state
    + creates an empty snapshot container
    # construction
  objdiff.add_target
    @ (state: diff_state, name: string, instructions: list[string]) -> diff_state
    + records the target (expected) instruction listing for a function
    # registration
  objdiff.add_current
    @ (state: diff_state, name: string, instructions: list[string]) -> diff_state
    + records the current (actual) instruction listing for a function
    # registration
  objdiff.compare_function
    @ (state: diff_state, name: string) -> result[function_diff, string]
    + returns a per-instruction diff with added, removed, and matched counts
    - returns error when the function is absent from either side
    # comparison
  objdiff.similarity
    @ (diff: function_diff) -> f64
    + returns a value in [0, 1] proportional to the fraction of matched instructions
    # scoring
  objdiff.summary
    @ (state: diff_state) -> list[tuple[string, f64]]
    + returns (function_name, similarity) pairs across all functions present on both sides
    + sorted by similarity ascending so the most divergent functions come first
    # reporting
