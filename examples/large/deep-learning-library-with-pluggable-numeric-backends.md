# Requirement: "a high-level deep learning library with multiple pluggable numeric backends"

A frontend for building and training neural networks; numeric operations are delegated to a backend interface.

std
  std.numeric
    std.numeric.zeros
      @ (shape: list[i32]) -> tensor
      + allocates a tensor filled with zeros
      # numeric
    std.numeric.randn
      @ (shape: list[i32], stddev: f32) -> tensor
      + allocates a tensor filled with gaussian noise
      # numeric
    std.numeric.matmul
      @ (a: tensor, b: tensor) -> tensor
      + returns the matrix product
      # numeric
    std.numeric.add
      @ (a: tensor, b: tensor) -> tensor
      + returns element-wise sum with broadcasting
      # numeric
    std.numeric.mul
      @ (a: tensor, b: tensor) -> tensor
      + returns element-wise product with broadcasting
      # numeric
    std.numeric.relu
      @ (x: tensor) -> tensor
      + applies rectified linear activation element-wise
      # numeric
    std.numeric.softmax
      @ (x: tensor, axis: i32) -> tensor
      + applies softmax along axis
      # numeric
    std.numeric.grad
      @ (loss: tensor, params: list[tensor]) -> list[tensor]
      + computes gradients of loss with respect to each parameter
      # autodiff

dl
  dl.backend_register
    @ (name: string, impl: backend_impl) -> void
    + registers a numeric backend under a name
    # backends
  dl.backend_select
    @ (name: string) -> result[void, string]
    + activates a previously registered backend
    - returns error when name is unknown
    # backends
  dl.dense_layer
    @ (input_dim: i32, output_dim: i32) -> layer
    + creates a fully connected layer with He-initialized weights and zero bias
    # layers
    -> std.numeric.randn
    -> std.numeric.zeros
  dl.activation_layer
    @ (kind: string) -> layer
    + creates an activation layer (one of "relu", "softmax")
    - returns an identity layer for unrecognized kinds
    # layers
    -> std.numeric.relu
    -> std.numeric.softmax
  dl.sequential
    @ (layers: list[layer]) -> model_state
    + builds a sequential model from a list of layers
    # model
  dl.forward
    @ (model: model_state, input: tensor) -> tensor
    + runs input through the model layers in order
    # model
    -> std.numeric.matmul
    -> std.numeric.add
  dl.loss_cross_entropy
    @ (logits: tensor, targets: tensor) -> tensor
    + returns the categorical cross-entropy loss
    # loss
  dl.loss_mse
    @ (predicted: tensor, targets: tensor) -> tensor
    + returns mean squared error
    # loss
  dl.optimizer_sgd
    @ (learning_rate: f32) -> optimizer_state
    + creates a stochastic gradient descent optimizer
    # optimization
  dl.optimizer_adam
    @ (learning_rate: f32, beta1: f32, beta2: f32) -> optimizer_state
    + creates an Adam optimizer with moment estimates
    # optimization
  dl.train_step
    @ (model: model_state, opt: optimizer_state, input: tensor, targets: tensor, loss_fn: fn(tensor, tensor) -> tensor) -> tuple[model_state, optimizer_state, f32]
    + runs a forward pass, computes gradients, applies the optimizer, and returns the loss value
    # training
    -> std.numeric.grad
    -> dl.forward
  dl.evaluate
    @ (model: model_state, input: tensor, targets: tensor) -> f32
    + returns the mean accuracy over a batch
    # evaluation
    -> dl.forward
