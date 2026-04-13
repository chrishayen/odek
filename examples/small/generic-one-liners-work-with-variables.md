# Requirement: "a utility library of small one-liners for common list and map operations"

A handful of generic collection helpers that would otherwise be rewritten per-call site.

std: (all units exist)

lang
  lang.map_list
    @ (xs: list[bytes], f: map_fn) -> list[bytes]
    + applies f to each element and returns the resulting list in order
    # lists
  lang.filter_list
    @ (xs: list[bytes], pred: pred_fn) -> list[bytes]
    + returns the elements for which pred returns true
    # lists
  lang.reduce_list
    @ (xs: list[bytes], init: bytes, step: reduce_fn) -> bytes
    + folds step over the list left-to-right starting from init
    # lists
  lang.contains
    @ (xs: list[bytes], needle: bytes) -> bool
    + returns true when any element equals needle
    # lists
  lang.map_keys
    @ (m: map[bytes, bytes]) -> list[bytes]
    + returns the keys of m with no ordering guarantee
    # maps
  lang.map_values
    @ (m: map[bytes, bytes]) -> list[bytes]
    + returns the values of m with no ordering guarantee
    # maps
  lang.default
    @ (value: optional[bytes], fallback: bytes) -> bytes
    + returns the contained value when present, otherwise fallback
    # optionals
