# Requirement: "a library for running large language models locally"

A local inference host: pulls and stores model weights, exposes a model registry, loads models into an inference backend, and offers a generation API. The backend (weights loader, tokenizer, transformer runner) is pluggable.

std
  std.fs
    std.fs.ensure_dir
      @ (path: string) -> result[void, string]
      + creates the directory if missing, including parents
      - returns error on permission failure
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file fully into memory
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the given path atomically
      # filesystem
  std.http
    std.http.download
      @ (url: string, dest: string, on_progress: fn(i64, i64) -> void) -> result[void, string]
      + streams the response body to disk, reporting bytes downloaded and total
      - returns error on non-2xx response
      # http
  std.hash
    std.hash.sha256_file
      @ (path: string) -> result[string, string]
      + returns the lowercase hex sha256 of the file's contents
      # hashing
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization

llmhost
  llmhost.new_store
    @ (root_dir: string) -> result[model_store, string]
    + opens or creates the on-disk model store at the given directory
    # storage
    -> std.fs.ensure_dir
  llmhost.index_lookup
    @ (store: model_store, name: string) -> optional[model_descriptor]
    + returns the descriptor for an installed model or none if not present
    # storage
  llmhost.pull
    @ (store: model_store, name: string, source_url: string, expected_sha256: string) -> result[model_descriptor, string]
    + downloads the model, verifies its digest, and registers it in the store
    - returns error when the digest does not match
    # acquisition
    -> std.http.download
    -> std.hash.sha256_file
    -> std.fs.write_all
  llmhost.load
    @ (store: model_store, name: string, backend: inference_backend) -> result[loaded_model, string]
    + materializes the model into the backend's runtime representation
    - returns error when the model is not in the store
    # loading
    -> std.fs.read_all
  llmhost.unload
    @ (model: loaded_model) -> void
    + releases the backend resources held by a loaded model
    # lifecycle
  llmhost.new_session
    @ (model: loaded_model, sampling: sampling_params) -> session_state
    + creates a generation session with sampling parameters and an empty kv cache
    # session
  llmhost.append_prompt
    @ (session: session_state, text: string) -> session_state
    + tokenizes and appends text to the session, extending the kv cache lazily
    # session
  llmhost.generate
    @ (session: session_state, max_new_tokens: i32) -> tuple[string, session_state]
    + runs the backend to produce up to max_new_tokens, returning the decoded text and advanced session
    + stops early at an end-of-stream token if the backend emits one
    # generation
  llmhost.stream
    @ (session: session_state, max_new_tokens: i32, on_token: fn(string) -> bool) -> session_state
    + like generate but invokes on_token for each decoded token; returning false from the callback halts generation
    # generation
  llmhost.parse_manifest
    @ (raw: string) -> result[model_descriptor, string]
    + parses a stored manifest into a model descriptor
    - returns error on malformed JSON
    # storage
    -> std.json.parse_object
