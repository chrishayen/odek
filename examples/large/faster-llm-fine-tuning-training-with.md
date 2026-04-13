# Requirement: "a memory-efficient fine-tuning trainer for large language models"

Loads a pretrained transformer, attaches low-rank adapter layers, streams a training corpus, runs gradient steps with mixed-precision accumulation, and writes updated adapter weights. Only the adapters are trained; base weights are held in a low-precision cache.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the file is missing
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path
      # filesystem
  std.tensor
    std.tensor.zeros
      @ (shape: list[i64], dtype: dtype) -> tensor
      + returns a tensor of the given shape filled with zeros
      # tensor
    std.tensor.matmul
      @ (a: tensor, b: tensor) -> result[tensor, string]
      + returns a @ b
      - returns error when inner dimensions do not match
      # tensor
    std.tensor.cast
      @ (t: tensor, target: dtype) -> tensor
      + returns a new tensor in the target precision
      # tensor

fine_tuner
  fine_tuner.load_base_model
    @ (weights_path: string) -> result[model_state, string]
    + loads pretrained weights into a low-precision cache
    - returns error when the weights file is malformed
    # model_loading
    -> std.fs.read_all
    -> std.tensor.cast
  fine_tuner.attach_adapters
    @ (model: model_state, rank: i32) -> model_state
    + inserts low-rank adapter matrices into every linear projection
    ? only adapter parameters are marked trainable
    # adapter_injection
    -> std.tensor.zeros
  fine_tuner.tokenize_corpus
    @ (corpus: string, max_len: i32) -> list[token_batch]
    + returns batches of packed token ids with attention masks
    # tokenization
  fine_tuner.forward
    @ (model: model_state, batch: token_batch) -> result[forward_output, string]
    + returns logits and per-token loss
    - returns error when batch shape is incompatible with the model
    # forward_pass
    -> std.tensor.matmul
  fine_tuner.backward
    @ (model: model_state, out: forward_output) -> gradients
    + returns gradients for adapter parameters only
    ? gradients are accumulated in fp32 even when forward runs in fp16
    # backward_pass
  fine_tuner.apply_step
    @ (model: model_state, grads: gradients, lr: f32) -> model_state
    + applies an optimizer step to the adapter parameters
    # optimizer_step
  fine_tuner.train
    @ (model: model_state, batches: list[token_batch], lr: f32, epochs: i32) -> result[training_report, string]
    + runs epochs over batches and returns per-step loss values
    - returns error when any forward pass fails
    # training_loop
  fine_tuner.save_adapters
    @ (model: model_state, path: string) -> result[void, string]
    + writes only the trained adapter weights to disk
    # checkpointing
    -> std.fs.write_all
  fine_tuner.load_adapters
    @ (model: model_state, path: string) -> result[model_state, string]
    + loads adapter weights into an already-initialized base model
    - returns error on shape mismatch
    # checkpointing
    -> std.fs.read_all
