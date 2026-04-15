# Requirement: "higher-order functions and operations on callable objects"

A small utility package for composing and memoizing callables.

std: (all units exist)

functools
  functools.compose
    fn (f: fn1, g: fn1) -> fn1
    + returns a function h such that h(x) = f(g(x))
    # composition
  functools.partial
    fn (f: fn_var, fixed: list[any_value]) -> fn_var
    + returns a function with the first len(fixed) arguments pre-bound
    # partial_application
  functools.memoize
    fn (f: fn1) -> fn1
    + returns a function that caches results keyed by input
    ? cache grows without bound; callers evict if needed
    # caching
  functools.reduce
    fn (f: reducer_fn, items: list[any_value], initial: any_value) -> any_value
    + returns the left fold of items under f starting from initial
    + returns initial for an empty list
    # reduction
  functools.lru_cache
    fn (f: fn1, capacity: i32) -> fn1
    + returns a function that caches up to capacity recent results, evicting least-recently-used
    - panics when capacity is non-positive
    # caching
