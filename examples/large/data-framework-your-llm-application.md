# Requirement: "a data framework for applications that use large language models"

Ingests documents, chunks them, embeds the chunks, stores them in a vector index, and retrieves the most relevant chunks for a query before synthesizing an answer through a model call.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file fully into memory
      - returns error when the path does not exist
      # filesystem
  std.text
    std.text.decode_utf8
      @ (data: bytes) -> result[string, string]
      + decodes bytes as utf-8
      - returns error on invalid utf-8
      # text
    std.text.split_whitespace
      @ (s: string) -> list[string]
      + splits a string into whitespace-separated tokens
      # text
  std.math
    std.math.dot_product
      @ (a: list[f32], b: list[f32]) -> f32
      + returns the dot product of two equal-length vectors
      - returns 0 when lengths differ
      # math
    std.math.vector_norm
      @ (v: list[f32]) -> f32
      + returns the Euclidean norm of the vector
      # math
  std.http
    std.http.post_json
      @ (url: string, headers: map[string, string], body: string) -> result[string, string]
      + sends a POST with a JSON body and returns the response body
      - returns error on non-2xx status
      # network

ragkit
  ragkit.load_text_file
    @ (path: string) -> result[document, string]
    + loads a plain text file as a document with the path as its id
    - returns error when the file cannot be read
    # ingest
    -> std.fs.read_all
    -> std.text.decode_utf8
  ragkit.chunk_document
    @ (doc: document, chunk_size: u32, overlap: u32) -> list[chunk]
    + splits the document into token-count-based chunks with the requested overlap
    + returns a single chunk when the document is smaller than chunk_size
    # chunking
    -> std.text.split_whitespace
  ragkit.embed_chunk
    @ (chunk: chunk, model_url: string, api_key: string) -> result[chunk, string]
    + attaches an embedding vector obtained from the remote model
    - returns error when the model endpoint rejects the request
    # embedding
    -> std.http.post_json
  ragkit.new_index
    @ () -> vector_index
    + creates an empty in-memory vector index
    # index
  ragkit.index_insert
    @ (index: vector_index, chunk: chunk) -> result[vector_index, string]
    + adds a chunk to the index
    - returns error when the chunk has no embedding
    # index
  ragkit.index_search
    @ (index: vector_index, query_vector: list[f32], top_k: u32) -> list[chunk]
    + returns the top-k chunks by cosine similarity
    + returns fewer than k when the index has fewer entries
    # retrieval
    -> std.math.dot_product
    -> std.math.vector_norm
  ragkit.build_prompt
    @ (question: string, context: list[chunk]) -> string
    + assembles a prompt with the question and retrieved context
    # prompting
  ragkit.answer
    @ (index: vector_index, question: string, model_url: string, api_key: string) -> result[string, string]
    + embeds the question, retrieves relevant chunks, and asks the model for an answer
    - returns error when the index is empty
    - returns error when the model call fails
    # pipeline
    -> std.http.post_json
