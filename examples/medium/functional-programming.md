# Requirement: "a functional programming primitives library"

Core higher-order list and optional operations. Each rune is a single generic primitive.

std: (all units exist)

fp
  fp.map
    @ (xs: list[T], f: func(T) -> U) -> list[U]
    + applies f to each element and returns the results in order
    + returns an empty list when xs is empty
    # transform
  fp.filter
    @ (xs: list[T], pred: func(T) -> bool) -> list[T]
    + returns elements for which pred returns true, preserving order
    # transform
  fp.fold_left
    @ (xs: list[T], initial: U, f: func(U, T) -> U) -> U
    + folds f over xs starting from initial, left to right
    + returns initial when xs is empty
    # reduce
  fp.compose
    @ (f: func(B) -> C, g: func(A) -> B) -> func(A) -> C
    + returns a function that applies g then f
    # composition
  fp.pipe
    @ (fs: list[func(T) -> T]) -> func(T) -> T
    + returns a function that applies each element of fs in order
    + returns identity when fs is empty
    # composition
  fp.map_optional
    @ (x: optional[T], f: func(T) -> U) -> optional[U]
    + returns some(f(v)) when x is some(v)
    - returns none when x is none
    # optional
  fp.flat_map
    @ (xs: list[T], f: func(T) -> list[U]) -> list[U]
    + concatenates the results of applying f to each element
    # transform
