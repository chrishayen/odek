# Requirement: "a library for fixed-size arrays whose length is part of the type"

Since the type system does not express length-parameterized generics, the project exposes a constructor that validates length and a bounded accessor.

std: (all units exist)

fixed_array
  fixed_array.new
    @ (length: i32, fill: i32) -> result[fixed_array_state, string]
    + creates an array of exactly length elements, all set to fill
    - returns error when length is negative
    # construction
  fixed_array.get
    @ (state: fixed_array_state, index: i32) -> result[i32, string]
    + returns the element at index
    - returns error when index is outside [0, length)
    # access
  fixed_array.set
    @ (state: fixed_array_state, index: i32, value: i32) -> result[fixed_array_state, string]
    + returns a new state with the element at index replaced
    - returns error when index is outside [0, length)
    # mutation
