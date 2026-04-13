# Requirement: "a graph database with pluggable storage backends"

A graph database stores (subject, predicate, object) triples and answers path queries. The storage backend is a narrow interface so in-memory, file, and key-value implementations can plug in.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire file as bytes
      - returns error when the path does not exist
      # io
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the file, truncating existing content
      # io
  std.encoding
    std.encoding.encode_triple
      @ (subject: string, predicate: string, object: string) -> bytes
      + encodes a triple as a length-prefixed byte sequence
      # encoding
    std.encoding.decode_triple
      @ (data: bytes) -> result[tuple[string, string, string], string]
      + decodes a length-prefixed triple
      - returns error on malformed input
      # encoding
  std.hash
    std.hash.fnv64
      @ (data: bytes) -> u64
      + computes a 64-bit FNV-1a hash
      # hashing

graph_db
  graph_db.open_memory
    @ () -> db_state
    + creates a new in-memory graph database
    # construction
  graph_db.open_file
    @ (path: string) -> result[db_state, string]
    + opens a graph database backed by a file at the given path
    - returns error when the file exists but is unreadable
    # construction
    -> std.fs.read_all
    -> std.encoding.decode_triple
  graph_db.add_triple
    @ (state: db_state, subject: string, predicate: string, object: string) -> db_state
    + inserts the triple; duplicates are idempotent
    # mutation
    -> std.hash.fnv64
  graph_db.remove_triple
    @ (state: db_state, subject: string, predicate: string, object: string) -> db_state
    + removes the matching triple
    - returns unchanged state when the triple is absent
    # mutation
  graph_db.query_out
    @ (state: db_state, subject: string, predicate: string) -> list[string]
    + returns all objects reachable from subject via the given predicate
    # query
  graph_db.query_in
    @ (state: db_state, predicate: string, object: string) -> list[string]
    + returns all subjects that reach object via the given predicate
    # query
  graph_db.traverse
    @ (state: db_state, start: string, predicate: string, max_depth: i32) -> list[string]
    + returns all nodes reachable from start by following predicate up to max_depth hops
    ? breadth-first; cycles are not revisited
    # traversal
  graph_db.count_triples
    @ (state: db_state) -> i64
    + returns the total number of triples currently stored
    # query
  graph_db.flush
    @ (state: db_state) -> result[void, string]
    + persists any pending writes to the backing store
    - returns error when the backing store rejects the write
    # persistence
    -> std.fs.write_all
    -> std.encoding.encode_triple
  graph_db.close
    @ (state: db_state) -> result[void, string]
    + flushes and releases all resources held by the database
    # lifecycle
