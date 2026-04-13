# Requirement: "a long-term memory store for agents keyed by user and embedding similarity"

Stores textual memories with vector embeddings per user and retrieves the most relevant items for a query.

std
  std.math
    std.math.sqrt
      @ (x: f64) -> f64
      + returns the non-negative square root
      # math
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

memory
  memory.new_store
    @ () -> memory_store
    + creates an empty store
    # construction
  memory.add
    @ (s: memory_store, user_id: string, text: string, embedding: list[f32]) -> tuple[memory_store, string]
    + stores the memory and returns a new memory id
    ? embedding length must match across all memories for a user
    # write
    -> std.time.now_seconds
  memory.get
    @ (s: memory_store, memory_id: string) -> optional[memory_record]
    + returns the stored record by id
    - returns none when unknown
    # read
  memory.delete
    @ (s: memory_store, memory_id: string) -> memory_store
    + removes the memory if present
    # write
  memory.cosine_similarity
    @ (a: list[f32], b: list[f32]) -> f32
    + returns dot(a,b) / (|a||b|)
    ? returns 0 when either vector has zero norm
    # similarity
    -> std.math.sqrt
  memory.search
    @ (s: memory_store, user_id: string, query_embedding: list[f32], top_k: i32) -> list[memory_record]
    + returns the top_k memories for the user ranked by cosine similarity
    # retrieval
  memory.forget_older_than
    @ (s: memory_store, user_id: string, max_age_seconds: i64) -> memory_store
    + drops memories for the user older than the given age
    # maintenance
    -> std.time.now_seconds
