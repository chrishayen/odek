# Requirement: "a flexible deep learning framework"

Expresses tensors, a static computation graph, and automatic differentiation. The std layer supplies tensor storage and linear-algebra primitives; the project layer composes them into layers, losses, and an optimizer loop.

std
  std.math
    std.math.exp
      @ (x: f64) -> f64
      + returns e to the x
      # math
    std.math.log
      @ (x: f64) -> f64
      + returns the natural log of x
      - returns NaN for x <= 0
      # math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the non-negative square root of x
      # math
  std.random
    std.random.new
      @ (seed: i64) -> rng_state
      + creates a deterministic RNG
      # random
    std.random.gaussian
      @ (state: rng_state) -> tuple[f64, rng_state]
      + returns a sample from the standard normal and the updated state
      # random
  std.tensor
    std.tensor.zeros
      @ (shape: list[i32]) -> tensor
      + returns a tensor filled with zeros
      # tensor
    std.tensor.from_values
      @ (shape: list[i32], values: list[f64]) -> result[tensor, string]
      + returns a tensor with the given row-major values
      - returns error when values length does not match shape
      # tensor
    std.tensor.matmul
      @ (a: tensor, b: tensor) -> result[tensor, string]
      + returns matrix multiplication a times b
      - returns error on incompatible shapes
      # tensor
    std.tensor.add
      @ (a: tensor, b: tensor) -> result[tensor, string]
      + elementwise sum with broadcasting
      - returns error on incompatible shapes
      # tensor
    std.tensor.map
      @ (a: tensor, f: f64_unary) -> tensor
      + applies f elementwise
      # tensor

dl
  dl.linear_layer
    @ (in_features: i32, out_features: i32, rng: rng_state) -> tuple[layer_params, rng_state]
    + initializes weights from a scaled normal and biases to zero
    # layers
    -> std.random.gaussian
    -> std.math.sqrt
  dl.linear_forward
    @ (params: layer_params, input: tensor) -> result[tensor, string]
    + returns input @ weight + bias
    # layers
    -> std.tensor.matmul
    -> std.tensor.add
  dl.relu
    @ (x: tensor) -> tensor
    + elementwise max(x, 0)
    # activations
    -> std.tensor.map
  dl.softmax
    @ (x: tensor) -> tensor
    + returns the softmax along the last axis
    ? stabilizes by subtracting the row max before exponentiating
    # activations
    -> std.math.exp
  dl.cross_entropy
    @ (logits: tensor, targets: list[i32]) -> f64
    + returns mean cross-entropy loss
    - returns +inf when a target index is out of range
    # loss
    -> std.math.log
  dl.build_graph
    @ (input_shape: list[i32], layers: list[layer_spec]) -> compute_graph
    + returns a graph that threads input_shape through each layer
    # graph
  dl.forward
    @ (graph: compute_graph, input: tensor, params: model_params) -> result[tensor, string]
    + runs the forward pass and returns the output tensor
    # evaluation
  dl.backward
    @ (graph: compute_graph, params: model_params, loss_grad: tensor) -> grad_map
    + returns gradients for every parameter via reverse-mode autodiff
    # autodiff
  dl.sgd_step
    @ (params: model_params, grads: grad_map, lr: f64) -> model_params
    + returns parameters with lr-scaled gradient subtracted
    # optimization
  dl.adam_step
    @ (state: adam_state, params: model_params, grads: grad_map, lr: f64) -> tuple[adam_state, model_params]
    + returns updated moments and parameters for one Adam step
    # optimization
    -> std.math.sqrt
  dl.fit_epoch
    @ (graph: compute_graph, params: model_params, data: batch_stream, lr: f64) -> tuple[model_params, f64]
    + runs one epoch and returns updated parameters and mean loss
    # training
