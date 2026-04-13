# Requirement: "a natural language processing library supporting latent semantic analysis"

LSA pipeline: tokenize, build a term-document matrix, compute TF-IDF, factor via SVD, and compare documents in the reduced space.

std
  std.math
    std.math.log
      @ (x: f64) -> f64
      + returns the natural logarithm
      # math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the square root
      # math
    std.math.dot
      @ (a: list[f64], b: list[f64]) -> f64
      + returns the dot product of two equal-length vectors
      - returns 0 when lengths differ
      # math
  std.linalg
    std.linalg.new_matrix
      @ (rows: i32, cols: i32) -> matrix
      + returns a zero-filled dense matrix
      # linear_algebra
    std.linalg.svd
      @ (m: matrix) -> svd_result
      + returns U, singular values, and V^T
      + singular values are in descending order
      # linear_algebra
    std.linalg.truncate_svd
      @ (svd: svd_result, k: i32) -> svd_result
      + keeps only the top k singular components
      # linear_algebra

nlp
  nlp.tokenize
    @ (text: string) -> list[string]
    + splits on whitespace and punctuation, lowercases, and drops empty tokens
    # tokenization
  nlp.remove_stopwords
    @ (tokens: list[string], stopwords: list[string]) -> list[string]
    + returns the tokens with any stopword removed
    # filtering
  nlp.build_vocab
    @ (documents: list[list[string]]) -> map[string,i32]
    + returns a deterministic token-to-index mapping
    # vocabulary
  nlp.term_document_matrix
    @ (documents: list[list[string]], vocab: map[string,i32]) -> matrix
    + fills a vocab_size-by-document_count matrix with raw term counts
    # representation
    -> std.linalg.new_matrix
  nlp.tfidf
    @ (counts: matrix) -> matrix
    + scales each cell by term frequency and inverse document frequency
    # weighting
    -> std.math.log
  nlp.fit_lsa
    @ (tfidf: matrix, k: i32) -> lsa_model
    + performs truncated SVD of the TF-IDF matrix to k dimensions
    - returns empty model when k <= 0
    # modeling
    -> std.linalg.svd
    -> std.linalg.truncate_svd
  nlp.project_document
    @ (model: lsa_model, counts: list[f64]) -> list[f64]
    + projects a raw count vector into the k-dimensional latent space
    # projection
  nlp.cosine_similarity
    @ (a: list[f64], b: list[f64]) -> f64
    + returns the cosine similarity of two vectors in [-1, 1]
    - returns 0 when either vector has zero magnitude
    # similarity
    -> std.math.dot
    -> std.math.sqrt
