# Requirement: "an n-dimensional array library with array views, multidimensional slicing, and efficient operations"

Dense row-major storage plus lightweight views that share the underlying buffer. Slicing returns a view; reshaping is view-only when the layout is contiguous.

std: (all units exist)

ndarray
  ndarray.new
    fn (shape: list[i32], fill: f64) -> ndarray_state
    + allocates a contiguous row-major buffer of product(shape) elements
    + every element is initialized to fill
    ? strides are computed as row-major from shape
    # construction
  ndarray.from_data
    fn (shape: list[i32], data: list[f64]) -> result[ndarray_state, string]
    + wraps an existing flat buffer
    - returns error when len(data) != product(shape)
    # construction
  ndarray.view
    fn (state: ndarray_state) -> ndarray_view
    + returns a full view over the array without copying
    # views
  ndarray.slice
    fn (view: ndarray_view, ranges: list[range]) -> result[ndarray_view, string]
    + returns a new view with offset and strides adjusted for the given per-axis ranges
    - returns error when a range is out of bounds or len(ranges) != rank
    # slicing
  ndarray.reshape
    fn (view: ndarray_view, new_shape: list[i32]) -> result[ndarray_view, string]
    + returns a view with the new shape when the current layout is contiguous
    - returns error when the element count would change
    - returns error when the layout is non-contiguous
    # views
  ndarray.get
    fn (view: ndarray_view, index: list[i32]) -> result[f64, string]
    + reads the element at index using the view's offset and strides
    - returns error when index is out of bounds
    # access
  ndarray.set
    fn (view: ndarray_view, index: list[i32], value: f64) -> result[void, string]
    + writes value at index in the underlying buffer
    - returns error when index is out of bounds
    # access
  ndarray.map
    fn (view: ndarray_view, f: unary_f64) -> ndarray_state
    + returns a freshly allocated contiguous array with f applied elementwise
    # elementwise
  ndarray.zip_add
    fn (a: ndarray_view, b: ndarray_view) -> result[ndarray_state, string]
    + returns a new contiguous array with elementwise sum
    - returns error when shapes disagree
    # elementwise
  ndarray.matmul
    fn (a: ndarray_view, b: ndarray_view) -> result[ndarray_state, string]
    + returns the matrix product when both views are rank 2 with compatible inner dims
    - returns error when either rank != 2 or inner dimensions disagree
    # linear_algebra
