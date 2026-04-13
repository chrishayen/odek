# Requirement: "a library that loads pre-trained diffusion models to generate and edit images, audio, and video"

High-level API for loading checkpoints by id and running conditional generation or edit pipelines across image, audio, and video modalities.

std
  std.fs
    std.fs.read_all_bytes
      @ (path: string) -> result[bytes, string]
      + returns the file contents as bytes
      - returns error when the file cannot be read
      # filesystem
    std.fs.write_all_bytes
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating it if needed
      - returns error on write failure
      # filesystem
  std.io
    std.io.http_download
      @ (url: string, dest_path: string) -> result[void, string]
      + streams a URL to a local file
      - returns error on network failure
      # http
  std.encoding
    std.encoding.tensor_deserialize
      @ (raw: bytes) -> result[tensor_map, string]
      + deserializes a tensor container into a name-to-tensor map
      - returns error on malformed data
      # serialization
  std.math
    std.math.randn
      @ (shape: list[i32], seed: i64) -> tensor
      + returns a tensor of standard-normal samples for the given shape
      # random

diffusion
  diffusion.load_checkpoint
    @ (model_id: string, cache_dir: string) -> result[model_handle, string]
    + downloads (if missing) and loads a checkpoint by id into a runnable handle
    - returns error when the model id is unknown or the download fails
    # loading
    -> std.io.http_download
    -> std.fs.read_all_bytes
    -> std.encoding.tensor_deserialize
  diffusion.make_scheduler
    @ (kind: string, num_steps: i32) -> result[scheduler_state, string]
    + creates a noise scheduler (ddpm, ddim, dpm) with the given step count
    - returns error when kind is unknown
    # scheduling
  diffusion.encode_prompt
    @ (model: model_handle, prompt: string) -> tensor
    + runs the text encoder and returns conditioning embeddings
    # conditioning
  diffusion.denoise_step
    @ (model: model_handle, scheduler: scheduler_state, latent: tensor, cond: tensor, step_index: i32) -> tensor
    + runs one reverse-diffusion step and returns the updated latent
    # inference
  diffusion.decode_latent_image
    @ (model: model_handle, latent: tensor) -> tensor
    + decodes an image latent into an RGB tensor
    # decoding
  diffusion.decode_latent_audio
    @ (model: model_handle, latent: tensor) -> tensor
    + decodes an audio latent into a waveform tensor
    # decoding
  diffusion.decode_latent_video
    @ (model: model_handle, latent: tensor) -> tensor
    + decodes a video latent into a stack of frames
    # decoding
  diffusion.generate_image
    @ (model: model_handle, prompt: string, steps: i32, seed: i64) -> result[tensor, string]
    + runs the full text-to-image pipeline and returns an RGB tensor
    - returns error when the model does not support image output
    # generation
    -> std.math.randn
  diffusion.generate_audio
    @ (model: model_handle, prompt: string, steps: i32, seed: i64) -> result[tensor, string]
    + runs the full text-to-audio pipeline and returns a waveform tensor
    - returns error when the model does not support audio output
    # generation
    -> std.math.randn
  diffusion.generate_video
    @ (model: model_handle, prompt: string, steps: i32, seed: i64) -> result[tensor, string]
    + runs the full text-to-video pipeline and returns a frame stack
    - returns error when the model does not support video output
    # generation
    -> std.math.randn
  diffusion.edit_image
    @ (model: model_handle, source: tensor, prompt: string, strength: f64, steps: i32, seed: i64) -> result[tensor, string]
    + runs an image-to-image pipeline conditioned on the source and prompt
    - returns error when strength is outside [0.0, 1.0]
    # editing
    -> std.math.randn
  diffusion.save_image_png
    @ (image: tensor, path: string) -> result[void, string]
    + writes an RGB tensor to disk as PNG
    - returns error on encoding or write failure
    # export
    -> std.fs.write_all_bytes
