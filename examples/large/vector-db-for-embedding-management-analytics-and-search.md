# Requirement: "a vector database for embedding management, analytics, and search"

A store for high-dimensional vectors with nearest-neighbor search. Persistence and the distance kernel are separated from index structure.

std
  std.math
    std.math.dot_product
      @ (a: list[f32], b: list[f32]) -> f32
      + returns sum of pairwise products
      - returns 0 when either list is empty
      # linear_algebra
    std.math.l2_norm
      @ (v: list[f32]) -> f32
      + returns sqrt of sum of squares
      # linear_algebra
    std.math.cosine_similarity
      @ (a: list[f32], b: list[f32]) -> f32
      + returns dot(a,b) / (norm(a)*norm(b))
      - returns 0 when either norm is 0
      # linear_algebra
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to path, creating or truncating
      # filesystem
  std.encoding
    std.encoding.encode_float_list
      @ (v: list[f32]) -> bytes
      + serializes a float list as length-prefixed little-endian f32 values
      # serialization
    std.encoding.decode_float_list
      @ (raw: bytes) -> result[list[f32], string]
      + decodes length-prefixed f32 list
      - returns error when the buffer is truncated
      # serialization

vector_db
  vector_db.new_collection
    @ (dim: i32, metric: string) -> result[collection_state, string]
    + creates an empty collection with the given dimension and metric
    - returns error when metric is not one of "cosine" or "l2"
    - returns error when dim <= 0
    # construction
  vector_db.insert
    @ (state: collection_state, id: string, vector: list[f32]) -> result[collection_state, string]
    + inserts id/vector and returns new state
    - returns error when vector length does not match collection dim
    - returns error when id already exists
    # insertion
  vector_db.delete
    @ (state: collection_state, id: string) -> result[collection_state, string]
    + removes the entry with the given id
    - returns error when id is not present
    # deletion
  vector_db.get
    @ (state: collection_state, id: string) -> optional[list[f32]]
    + returns stored vector when present
    - returns empty when id is unknown
    # lookup
  vector_db.search
    @ (state: collection_state, query: list[f32], k: i32) -> result[list[tuple[string, f32]], string]
    + returns up to k (id, score) pairs ordered best-first
    - returns error when query length does not match dim
    ? linear scan is acceptable; indexing is out of scope
    # nearest_neighbor_search
    -> std.math.cosine_similarity
    -> std.math.dot_product
    -> std.math.l2_norm
  vector_db.count
    @ (state: collection_state) -> i64
    + returns number of vectors in the collection
    # analytics
  vector_db.mean_vector
    @ (state: collection_state) -> optional[list[f32]]
    + returns the component-wise mean of all stored vectors
    - returns empty when collection is empty
    # analytics
  vector_db.save
    @ (state: collection_state, path: string) -> result[void, string]
    + serializes ids and vectors to the given path
    # persistence
    -> std.encoding.encode_float_list
    -> std.fs.write_all
  vector_db.load
    @ (path: string) -> result[collection_state, string]
    + reconstructs a collection from a saved file
    - returns error when the file is corrupt
    # persistence
    -> std.fs.read_all
    -> std.encoding.decode_float_list
