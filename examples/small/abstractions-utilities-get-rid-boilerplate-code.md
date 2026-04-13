# Requirement: "a small utilities library that reduces boilerplate in typical business logic"

A restrained set of genuinely generic helpers — optional unwrapping with a default, result chaining, and a list partitioner. Anything more specific belongs in the caller's code, not here.

std: (all units exist)

boilerless
  boilerless.or_default
    @ (opt: optional[string], fallback: string) -> string
    + returns the wrapped value when present
    + returns the fallback when the option is empty
    # option
  boilerless.first_error
    @ (results: list[result[void, string]]) -> optional[string]
    + returns the first error message found in the list
    - returns none when every result is ok
    # result
  boilerless.partition
    @ (items: list[string], predicate: fn(string) -> bool) -> tuple[list[string], list[string]]
    + returns (matching, non_matching) preserving input order in both lists
    # collection
  boilerless.group_by
    @ (items: list[string], key_fn: fn(string) -> string) -> map[string, list[string]]
    + groups items by the key produced by key_fn
    + preserves input order within each group
    # collection
