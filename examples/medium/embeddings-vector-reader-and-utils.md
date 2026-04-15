# Requirement: "a reader and utility functions for word embedding vectors"

Loads a binary embeddings file into an in-memory table and supports lookup plus cosine similarity operations.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + returns the file contents as bytes
      - returns error when the file cannot be read
      # filesystem

embeddings
  embeddings.load
    fn (path: string) -> result[embedding_table, string]
    + parses a binary embeddings file into a vocabulary and vector matrix
    - returns error when the header cannot be read
    - returns error when a vector row is truncated
    # loading
    -> std.fs.read_all
  embeddings.vector_for
    fn (table: embedding_table, word: string) -> optional[list[f32]]
    + returns the vector for a word
    - returns none when the word is not in the vocabulary
    # lookup
  embeddings.cosine_similarity
    fn (a: list[f32], b: list[f32]) -> result[f32, string]
    + returns the cosine similarity of two vectors
    - returns error when the vectors have different dimensions
    - returns error when either vector has zero magnitude
    # similarity
  embeddings.most_similar
    fn (table: embedding_table, word: string, k: i32) -> result[list[tuple[string, f32]], string]
    + returns the top-k words by cosine similarity, excluding the query word
    - returns error when the query word is not in the vocabulary
    # similarity
