# Requirement: "a deep learning framework that organizes training, evaluation, and inference for differentiable models"

The project layer owns the training loop, metric aggregation, checkpointing, and inference. Core tensor ops live in std.

std
  std.tensor
    std.tensor.zeros
      @ (shape: list[i32]) -> tensor_state
      + returns a dense tensor of the given shape filled with zero
      # tensor
    std.tensor.matmul
      @ (a: tensor_state, b: tensor_state) -> result[tensor_state, string]
      + returns the matrix product of two 2-D tensors
      - returns error when inner dimensions do not match
      # tensor
    std.tensor.add
      @ (a: tensor_state, b: tensor_state) -> result[tensor_state, string]
      + returns the elementwise sum of two tensors
      - returns error when shapes do not match
      # tensor
    std.tensor.relu
      @ (x: tensor_state) -> tensor_state
      + returns elementwise max(x, 0)
      # tensor
    std.tensor.softmax_cross_entropy
      @ (logits: tensor_state, labels: list[i32]) -> f64
      + returns mean cross-entropy loss
      # tensor
  std.autograd
    std.autograd.backward
      @ (loss: tensor_state) -> grad_state
      + computes gradients for all leaf tensors contributing to loss
      # autograd
    std.autograd.apply_sgd
      @ (grads: grad_state, learning_rate: f64) -> void
      + applies an SGD step to parameters in-place using accumulated gradients
      # optimization
  std.fs
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, replacing it if it exists
      - returns error on I/O failure
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads an entire file as bytes
      - returns error when the file does not exist
      # filesystem

dl_framework
  dl_framework.new_model
    @ () -> model_state
    + creates an empty sequential model
    # construction
  dl_framework.add_dense
    @ (model: model_state, in_dim: i32, out_dim: i32) -> model_state
    + appends a dense layer with random initialization
    # architecture
  dl_framework.add_relu
    @ (model: model_state) -> model_state
    + appends a ReLU activation layer
    # architecture
    -> std.tensor.relu
  dl_framework.forward
    @ (model: model_state, input: tensor_state) -> result[tensor_state, string]
    + runs the model forward and returns its output tensor
    - returns error when input shape does not match the first layer
    # inference
    -> std.tensor.matmul
    -> std.tensor.add
  dl_framework.train_step
    @ (model: model_state, batch: tensor_state, labels: list[i32], learning_rate: f64) -> result[f64, string]
    + runs forward, computes loss, back-propagates, and applies SGD
    + returns the scalar loss for the batch
    - returns error on any shape mismatch along the path
    # training
    -> std.tensor.softmax_cross_entropy
    -> std.autograd.backward
    -> std.autograd.apply_sgd
  dl_framework.evaluate
    @ (model: model_state, batch: tensor_state, labels: list[i32]) -> result[f64, string]
    + returns the top-1 accuracy on a batch
    - returns error on shape mismatch
    # evaluation
  dl_framework.new_dataset
    @ () -> dataset_state
    + creates an empty dataset
    # construction
  dl_framework.add_example
    @ (dataset: dataset_state, features: list[f64], label: i32) -> dataset_state
    + appends one (features, label) example
    # dataset
  dl_framework.iterate_batches
    @ (dataset: dataset_state, batch_size: i32) -> list[tuple[tensor_state, list[i32]]]
    + returns successive batches of (features tensor, labels list)
    ? the final batch may be smaller than batch_size
    # dataset
  dl_framework.save_checkpoint
    @ (model: model_state, path: string) -> result[void, string]
    + writes model parameters to a file
    - returns error on I/O failure
    # persistence
    -> std.fs.write_all
  dl_framework.load_checkpoint
    @ (path: string) -> result[model_state, string]
    + reads model parameters from a file
    - returns error when the file is missing or corrupt
    # persistence
    -> std.fs.read_all
  dl_framework.predict
    @ (model: model_state, features: list[f64]) -> result[i32, string]
    + returns the argmax class for a single input
    - returns error when feature dimension is wrong
    # inference
