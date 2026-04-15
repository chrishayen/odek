# Requirement: "a framework for loading and running pre-trained models for text, vision, and audio tasks"

A model hub and inference dispatcher. Models are registered by name with a modality tag; inputs are preprocessed, fed to an opaque model handle, and the output is decoded into a task-appropriate result.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file entirely into memory
      - returns error when the file does not exist or is unreadable
      # filesystem
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string map
      - returns error on malformed JSON
      # serialization
  std.tensor
    std.tensor.zeros
      fn (shape: list[i32]) -> tensor
      + returns a tensor of the given shape filled with zeros
      # tensors
    std.tensor.from_floats
      fn (shape: list[i32], values: list[f32]) -> result[tensor, string]
      + builds a tensor from a flat list of floats
      - returns error when values.length does not match product(shape)
      # tensors
    std.tensor.argmax
      fn (t: tensor, axis: i32) -> list[i32]
      + returns the index of the maximum value along axis for each slice
      # tensors

models
  models.new_hub
    fn () -> hub_state
    + creates an empty hub with no registered models
    # construction
  models.register
    fn (state: hub_state, name: string, modality: string, config_path: string, weights_path: string) -> hub_state
    + records a model's metadata and on-disk locations
    - returns unchanged state when name is already registered
    # registry
  models.load
    fn (state: hub_state, name: string) -> result[model_handle, string]
    + reads config and weights from disk, returning an opaque handle
    - returns error when name is not registered
    - returns error when either file cannot be read
    # loading
    -> std.fs.read_all
    -> std.json.parse_object
  models.unload
    fn (state: hub_state, handle: model_handle) -> hub_state
    + releases resources held by a loaded model
    # loading
  models.tokenize_text
    fn (handle: model_handle, text: string) -> list[i32]
    + converts text into token ids using the model's tokenizer
    # text_preprocessing
  models.detokenize
    fn (handle: model_handle, ids: list[i32]) -> string
    + reconstructs text from token ids
    # text_postprocessing
  models.preprocess_image
    fn (handle: model_handle, pixels: bytes, width: i32, height: i32) -> result[tensor, string]
    + resizes and normalizes an image into the model's expected input tensor
    - returns error when pixels length does not match width*height*3
    # vision_preprocessing
    -> std.tensor.from_floats
  models.preprocess_audio
    fn (handle: model_handle, samples: list[f32], sample_rate: i32) -> result[tensor, string]
    + resamples and frames audio into the model's input tensor
    - returns error when sample_rate is non-positive
    # audio_preprocessing
    -> std.tensor.from_floats
  models.run_inference
    fn (handle: model_handle, input: tensor) -> result[tensor, string]
    + returns the model's raw output tensor
    - returns error when input shape does not match the model's expected input
    # inference
  models.classify
    fn (handle: model_handle, logits: tensor) -> result[list[string], string]
    + turns logits into a ranked list of label strings
    - returns error when the model has no label vocabulary
    # classification
    -> std.tensor.argmax
  models.generate_text
    fn (handle: model_handle, prompt_ids: list[i32], max_new_tokens: i32) -> result[list[i32], string]
    + autoregressively produces up to max_new_tokens new token ids
    - returns error when max_new_tokens is not positive
    # text_generation
