# Requirement: "a serving framework for large language models and multimodal inputs"

Continuous-batching inference server. The model runtime itself is an opaque backend; the library coordinates queuing, batching, and streaming.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.sync
    std.sync.channel_new
      @ () -> channel_handle
      + creates an unbounded channel carrying opaque messages
      # concurrency
    std.sync.channel_send
      @ (ch: channel_handle, msg: bytes) -> void
      + enqueues a message
      # concurrency
    std.sync.channel_recv
      @ (ch: channel_handle) -> optional[bytes]
      + returns the next message or none if the channel is closed
      # concurrency

serving
  serving.new_engine
    @ (backend: model_backend, max_batch: i32, max_seq_len: i32) -> engine_state
    + creates an engine with an empty queue and the given limits
    # construction
  serving.enqueue_text_request
    @ (eng: engine_state, prompt: string, max_tokens: i32) -> request_handle
    + admits a text generation request and returns a handle for streaming output
    # queuing
  serving.enqueue_multimodal_request
    @ (eng: engine_state, prompt: string, images: list[bytes], max_tokens: i32) -> request_handle
    + admits a request carrying encoded images alongside the text prompt
    # queuing
  serving.tokenize_prompt
    @ (eng: engine_state, prompt: string) -> list[i32]
    + delegates to the backend tokenizer
    # tokenization
  serving.form_batch
    @ (eng: engine_state) -> batch
    + selects up to max_batch pending requests that fit in the token budget
    + pads sequences to the longest in the batch
    # batching
  serving.step_batch
    @ (eng: engine_state, b: batch) -> list[token_step]
    + advances the batch by one decode step and returns the new token for each active request
    - returns empty when the batch has no active requests
    # inference
  serving.stream_output
    @ (eng: engine_state, h: request_handle) -> optional[string]
    + returns the next decoded token string, or none when the request has finished
    # streaming
    -> std.sync.channel_recv
  serving.cancel
    @ (eng: engine_state, h: request_handle) -> void
    + marks a request as cancelled so the next batch drops it
    # control
  serving.run_loop
    @ (eng: engine_state) -> void
    + repeatedly forms batches, steps them, and fans tokens out to streaming channels
    # server
    -> serving.form_batch
    -> serving.step_batch
    -> std.sync.channel_send
    -> std.time.now_millis
