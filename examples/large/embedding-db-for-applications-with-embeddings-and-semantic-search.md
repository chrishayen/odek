# Requirement: "an embedding database for building applications with embeddings and semantic search"

A vector store with collections, upsert, and k-nearest-neighbor search. std provides math and storage primitives.

std
  std.math
    std.math.dot_product
      @ (a: list[f32], b: list[f32]) -> f32
      + returns sum of pairwise products
      - returns 0 when lengths mismatch
      # math
    std.math.l2_norm
      @ (v: list[f32]) -> f32
      + returns euclidean length of the vector
      # math
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error if file is missing
      # io
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes atomically
      # io
  std.encoding
    std.encoding.encode_f32_list
      @ (v: list[f32]) -> bytes
      + little-endian packed floats
      # encoding
    std.encoding.decode_f32_list
      @ (data: bytes) -> result[list[f32], string]
      + decodes little-endian packed floats
      - returns error on non-multiple-of-4 length
      # encoding

embedding_db
  embedding_db.create_collection
    @ (name: string, dim: i32) -> result[collection_handle, string]
    + creates an empty collection with the given vector dimension
    - returns error when dim is not positive
    # collection_management
  embedding_db.upsert
    @ (col: collection_handle, id: string, vector: list[f32], metadata: map[string, string]) -> result[void, string]
    + inserts or replaces an entry by id
    - returns error when vector length does not match collection dim
    # write_path
  embedding_db.delete
    @ (col: collection_handle, id: string) -> result[bool, string]
    + returns true when an entry was removed
    + returns false when id was not present
    # write_path
  embedding_db.cosine_similarity
    @ (a: list[f32], b: list[f32]) -> f32
    + returns dot(a,b) / (||a|| * ||b||)
    - returns 0 when either vector has zero norm
    # similarity
    -> std.math.dot_product
    -> std.math.l2_norm
  embedding_db.query_nearest
    @ (col: collection_handle, vector: list[f32], k: i32) -> result[list[query_hit], string]
    + returns top-k hits sorted by descending similarity
    - returns error when vector length does not match dim
    ? brute-force scan; no ANN index
    # search
    -> std.math.dot_product
    -> std.math.l2_norm
  embedding_db.filter_by_metadata
    @ (hits: list[query_hit], key: string, value: string) -> list[query_hit]
    + returns only hits whose metadata[key] equals value
    # filtering
  embedding_db.save_to_disk
    @ (col: collection_handle, path: string) -> result[void, string]
    + persists the collection to a single file
    - returns error on io failure
    # persistence
    -> std.encoding.encode_f32_list
    -> std.fs.write_all
  embedding_db.load_from_disk
    @ (path: string) -> result[collection_handle, string]
    + reconstructs a collection from a saved file
    - returns error when the file is corrupt or incompatible
    # persistence
    -> std.fs.read_all
    -> std.encoding.decode_f32_list
