# Requirement: "an embedded document-graph database"

An in-process database that stores documents keyed by table and id, supports graph edges between them, and runs simple filter queries. Persistence and serialization go through std primitives.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of the file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes the bytes atomically
      - returns error when the parent directory does not exist
      # filesystem
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses any JSON value into an opaque tree
      - returns error on invalid JSON
      # serialization
    std.json.encode_value
      @ (value: json_value) -> string
      + serializes a JSON value to canonical text
      # serialization
  std.hash
    std.hash.fnv64
      @ (data: bytes) -> u64
      + computes the 64-bit FNV-1a hash
      # hashing

docgraph
  docgraph.open
    @ (path: string) -> result[db_state, string]
    + creates or loads a database at the given path
    - returns error when the file is corrupted
    # storage
    -> std.fs.read_all
    -> std.json.parse_value
  docgraph.close
    @ (db: db_state) -> result[void, string]
    + flushes and releases the database
    # lifecycle
    -> std.fs.write_all
    -> std.json.encode_value
  docgraph.create
    @ (db: db_state, table: string, doc: json_value) -> result[tuple[string, db_state], string]
    + inserts the document and returns the generated id
    - returns error when the table name contains invalid characters
    # mutation
    -> std.hash.fnv64
  docgraph.get
    @ (db: db_state, table: string, id: string) -> result[json_value, string]
    + returns the document with the matching table and id
    - returns error when no such document exists
    # query
  docgraph.update
    @ (db: db_state, table: string, id: string, doc: json_value) -> result[db_state, string]
    + replaces the document in place
    - returns error when the document does not exist
    # mutation
  docgraph.delete
    @ (db: db_state, table: string, id: string) -> result[db_state, string]
    + removes the document and any edges touching it
    - returns error when the document does not exist
    # mutation
  docgraph.relate
    @ (db: db_state, from_table: string, from_id: string, edge: string, to_table: string, to_id: string) -> result[db_state, string]
    + creates a directed edge labeled by edge between the two documents
    - returns error when either endpoint does not exist
    # graph
  docgraph.neighbors
    @ (db: db_state, table: string, id: string, edge: string) -> result[list[json_value], string]
    + returns target documents reachable via the labeled edge
    - returns error when the source document does not exist
    # graph
  docgraph.select_where
    @ (db: db_state, table: string, field: string, equals: json_value) -> list[json_value]
    + returns all documents in the table whose field matches the value
    # query
  docgraph.count
    @ (db: db_state, table: string) -> i64
    + returns the number of documents in the table
    # query
