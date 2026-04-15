# Requirement: "a graph-plus-vector database for storing nodes, edges, and embeddings for similarity search"

Project layer stores nodes, edges, and vector embeddings and supports similarity queries over stored vectors. Low-level math primitives live in std.

std
  std.math
    std.math.dot
      fn (a: list[f32], b: list[f32]) -> f32
      + returns the dot product of two equal-length vectors
      - returns 0 when either input is empty
      # math
    std.math.norm
      fn (v: list[f32]) -> f32
      + returns the L2 norm of a vector
      # math

graphvec
  graphvec.new
    fn () -> db_state
    + returns an empty database
    # construction
  graphvec.add_node
    fn (db: db_state, id: string, props: map[string, string], embedding: list[f32]) -> db_state
    + stores a node with its properties and associated embedding
    ? an empty embedding means the node is not indexed for similarity
    # write
  graphvec.add_edge
    fn (db: db_state, source_id: string, target_id: string, label: string) -> result[db_state, string]
    + adds a labelled directed edge between two existing nodes
    - returns error when either endpoint does not exist
    # write
  graphvec.get_node
    fn (db: db_state, id: string) -> optional[map[string, string]]
    + returns the node's properties when it exists
    - returns none for unknown ids
    # read
  graphvec.neighbors
    fn (db: db_state, id: string, label: string) -> list[string]
    + returns the target ids of edges out of id with the given label
    # read
  graphvec.cosine
    fn (a: list[f32], b: list[f32]) -> f32
    + returns cosine similarity between two vectors
    - returns 0 when either vector has zero norm
    # similarity
    -> std.math.dot
    -> std.math.norm
  graphvec.search_similar
    fn (db: db_state, query: list[f32], k: i32) -> list[tuple[string, f32]]
    + returns the top-k node ids ranked by cosine similarity to the query
    ? ties broken by insertion order
    # similarity
  graphvec.search_similar_filtered
    fn (db: db_state, query: list[f32], k: i32, required_prop: string, required_value: string) -> list[tuple[string, f32]]
    + returns the top-k similar nodes whose property matches the given key-value pair
    # similarity
