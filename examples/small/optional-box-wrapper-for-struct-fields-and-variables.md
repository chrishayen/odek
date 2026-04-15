# Requirement: "an optional value wrapper for struct fields and variables"

A tiny wrapper that represents "present or absent" and the operations callers actually need.

std: (all units exist)

optional_box
  optional_box.some
    fn (value: string) -> optional[string]
    + wraps a value as a present optional
    # construction
  optional_box.none
    fn () -> optional[string]
    + returns an empty optional
    # construction
  optional_box.unwrap_or
    fn (opt: optional[string], fallback: string) -> string
    + returns the inner value when present
    - returns fallback when the optional is empty
    # access
  optional_box.map
    fn (opt: optional[string], transform: function[string, string]) -> optional[string]
    + applies transform when present
    - returns an empty optional when input is empty
    # transformation
