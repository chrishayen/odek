# Requirement: "a numerical computing library with automatic differentiation"

Dense n-dimensional arrays, elementwise and reduction operations, and reverse-mode automatic differentiation over a computation graph. JIT compilation is out of scope for this library level — the caller can lower the recorded graph elsewhere.

std: (all units exist)

numeric
  numeric.tensor_zeros
    fn (shape: list[i64]) -> tensor
    + creates a tensor filled with zeros of the given shape
    # construction
  numeric.tensor_from_values
    fn (shape: list[i64], values: list[f64]) -> result[tensor, string]
    + creates a tensor with the given shape and row-major values
    - returns error when product of shape differs from values length
    # construction
  numeric.tensor_shape
    fn (t: tensor) -> list[i64]
    + returns the tensor shape
    # inspection
  numeric.add
    fn (a: tensor, b: tensor) -> result[tensor, string]
    + returns elementwise sum when shapes match or broadcast
    - returns error on incompatible shapes
    # arithmetic
  numeric.mul
    fn (a: tensor, b: tensor) -> result[tensor, string]
    + returns elementwise product when shapes match or broadcast
    - returns error on incompatible shapes
    # arithmetic
  numeric.matmul
    fn (a: tensor, b: tensor) -> result[tensor, string]
    + returns matrix product of two 2-D tensors
    - returns error when inner dimensions disagree
    # linear_algebra
  numeric.sum
    fn (t: tensor, axis: i32) -> tensor
    + returns the sum along the given axis
    # reduction
  numeric.relu
    fn (t: tensor) -> tensor
    + returns elementwise max(x, 0)
    # activation
  numeric.new_graph
    fn () -> graph_state
    + creates an empty computation graph
    # autodiff
  numeric.param
    fn (graph: graph_state, initial: tensor) -> tuple[node_ref, graph_state]
    + registers a trainable parameter node
    # autodiff
  numeric.constant
    fn (graph: graph_state, value: tensor) -> tuple[node_ref, graph_state]
    + registers a non-differentiable constant node
    # autodiff
  numeric.op
    fn (graph: graph_state, op_name: string, inputs: list[node_ref]) -> tuple[node_ref, graph_state]
    + records an op node with the given inputs
    - returns unchanged state on unknown op name
    # autodiff
  numeric.forward
    fn (graph: graph_state, target: node_ref) -> result[tensor, string]
    + evaluates the graph up to target using recorded ops
    - returns error when a required input is missing
    # autodiff
  numeric.backward
    fn (graph: graph_state, loss: node_ref) -> map[string, tensor]
    + returns gradients of loss with respect to each parameter
    ? uses reverse-mode accumulation over the topologically sorted graph
    # autodiff
