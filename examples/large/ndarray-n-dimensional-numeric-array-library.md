# Requirement: "an n-dimensional numeric array library"

A dense array type with shape, elementwise arithmetic, reductions, and basic linear algebra. Element type is f64.

std: (all units exist)

ndarray
  ndarray.from_values
    @ (shape: list[i32], values: list[f64]) -> result[array, string]
    + returns an array whose row-major contents match values
    - returns error when values length does not match the product of shape
    # construction
  ndarray.zeros
    @ (shape: list[i32]) -> array
    + returns an array of the given shape filled with 0.0
    # construction
  ndarray.ones
    @ (shape: list[i32]) -> array
    + returns an array of the given shape filled with 1.0
    # construction
  ndarray.shape_of
    @ (a: array) -> list[i32]
    + returns the array's dimensions
    # query
  ndarray.get
    @ (a: array, indices: list[i32]) -> result[f64, string]
    + returns the element at the given row-major index
    - returns error when indices length or any value is out of bounds
    # access
  ndarray.set
    @ (a: array, indices: list[i32], value: f64) -> result[array, string]
    + returns a new array with the element at indices replaced
    - returns error when indices are out of bounds
    # access
  ndarray.reshape
    @ (a: array, new_shape: list[i32]) -> result[array, string]
    + returns a new view with the requested shape
    - returns error when the total element count changes
    # transform
  ndarray.add
    @ (left: array, right: array) -> result[array, string]
    + returns the elementwise sum of two arrays with the same shape
    - returns error when shapes do not match
    # arithmetic
  ndarray.mul
    @ (left: array, right: array) -> result[array, string]
    + returns the elementwise product of two arrays with the same shape
    - returns error when shapes do not match
    # arithmetic
  ndarray.scalar_mul
    @ (a: array, scalar: f64) -> array
    + returns an array with every element multiplied by scalar
    # arithmetic
  ndarray.sum
    @ (a: array) -> f64
    + returns the sum of all elements
    + returns 0.0 for a zero-element array
    # reduction
  ndarray.mean
    @ (a: array) -> result[f64, string]
    + returns the arithmetic mean of all elements
    - returns error when the array has zero elements
    # reduction
  ndarray.matmul
    @ (left: array, right: array) -> result[array, string]
    + returns the matrix product for two 2-D arrays
    - returns error when either input is not 2-D
    - returns error when the inner dimensions do not match
    # linear_algebra
  ndarray.transpose
    @ (a: array) -> result[array, string]
    + returns the transpose of a 2-D array
    - returns error when a is not 2-D
    # linear_algebra
