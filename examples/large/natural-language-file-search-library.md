# Requirement: "a natural-language file search library"

Indexes a directory tree by embedding file contents, then answers free-form queries by ranking files via vector similarity.

std
  std.fs
    std.fs.walk
      fn (root: string) -> result[list[string], string]
      + yields all regular file paths under root, recursively
      - returns error when root is not a directory
      # filesystem
    std.fs.read_text
      fn (path: string) -> result[string, string]
      + reads file as UTF-8 text
      - returns error on non-text or missing file
      # filesystem
  std.text
    std.text.chunk
      fn (text: string, max_chars: i32, overlap: i32) -> list[string]
      + splits text into overlapping chunks of at most max_chars
      + returns one chunk when text is shorter than max_chars
      # text_processing
  std.math
    std.math.cosine_similarity
      fn (a: list[f32], b: list[f32]) -> f32
      + returns cosine similarity in [-1.0, 1.0]
      - returns 0.0 when either vector is all zeros
      # math

file_search
  file_search.new_index
    fn (embed_dim: i32) -> index_state
    + creates an empty index with the given embedding dimension
    # construction
  file_search.add_document
    fn (state: index_state, path: string, chunks: list[string], embeddings: list[list[f32]]) -> result[index_state, string]
    + stores each chunk with its embedding under the given file path
    - returns error when chunk count and embedding count differ
    - returns error when any embedding does not match embed_dim
    # indexing
  file_search.index_directory
    fn (state: index_state, root: string, embed: fn(string) -> list[f32]) -> result[index_state, string]
    + walks root, chunks each file, computes embeddings, adds all chunks
    - returns error when root cannot be read
    ? caller supplies the embedding function so the library is model-agnostic
    # indexing
    -> std.fs.walk
    -> std.fs.read_text
    -> std.text.chunk
  file_search.query
    fn (state: index_state, query_embedding: list[f32], top_k: i32) -> list[search_hit]
    + returns the top_k chunks ranked by cosine similarity, descending
    + returns fewer than top_k when the index has fewer chunks
    - returns empty list when the index is empty
    # retrieval
    -> std.math.cosine_similarity
  file_search.top_files
    fn (hits: list[search_hit]) -> list[string]
    + collapses chunk hits to unique file paths preserving rank order
    # retrieval
