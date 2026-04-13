# Requirement: "an approximate nearest neighbor index optimized for low memory use"

Builds a forest of random projection trees over high-dimensional vectors and answers nearest-neighbor queries by descending each tree and merging candidates.

std
  std.math
    std.math.dot
      @ (a: list[f32], b: list[f32]) -> f32
      + returns the dot product of two equal-length vectors
      # math
    std.math.l2_distance
      @ (a: list[f32], b: list[f32]) -> f32
      + returns the Euclidean distance between two equal-length vectors
      # math
  std.random
    std.random.new
      @ (seed: i64) -> rng_state
      + creates a deterministic random generator from a seed
      # random
    std.random.uniform
      @ (rng: rng_state) -> tuple[f32, rng_state]
      + returns a float in [0,1) and the advanced state
      # random
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of the file
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes the full contents to the file
      # filesystem

annindex
  annindex.new
    @ (dim: i32, metric: string) -> result[ann_builder, string]
    + creates an empty builder for vectors of the given dimension
    - returns error when metric is not one of "l2" or "angular"
    # construction
  annindex.add_item
    @ (builder: ann_builder, id: i64, vector: list[f32]) -> result[ann_builder, string]
    + adds one vector to the builder
    - returns error when the vector length does not match dim
    # ingestion
  annindex.build
    @ (builder: ann_builder, num_trees: i32, seed: i64) -> result[ann_index, string]
    + constructs the projection-tree forest and returns a queryable index
    - returns error when the builder is empty
    # build
    -> std.random.new
    -> std.random.uniform
  annindex.random_split_plane
    @ (rng: rng_state, vectors: list[list[f32]]) -> tuple[split_plane, rng_state]
    + picks two random vectors from the candidate set and returns the midplane between them
    # build
    -> std.math.dot
  annindex.descend_tree
    @ (tree: ann_tree, query: list[f32], k: i32) -> list[i64]
    + returns up to k candidate ids from a single tree by descending to the leaf
    # query
    -> std.math.dot
  annindex.query
    @ (index: ann_index, vector: list[f32], k: i32) -> result[list[i64], string]
    + returns the k nearest ids, merging candidates from every tree and re-ranking by true distance
    - returns error when the vector length does not match the index dimension
    # query
    -> std.math.l2_distance
  annindex.save
    @ (index: ann_index, path: string) -> result[void, string]
    + serializes the index to a file
    # persistence
    -> std.fs.write_all
  annindex.load
    @ (path: string) -> result[ann_index, string]
    + reads a previously saved index from a file
    - returns error on corrupt file
    # persistence
    -> std.fs.read_all
