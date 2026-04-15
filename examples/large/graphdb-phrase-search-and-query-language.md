# Requirement: "a transactional graph database with phrase search and a query language"

Stores nodes and edges, supports atomic transactions, indexes text for phrase search, and exposes a small query language.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads a file fully into memory
      - returns error on missing file
      # filesystem
    std.fs.append
      fn (path: string, data: bytes) -> result[void, string]
      + appends bytes to a file, creating it if missing
      # filesystem
    std.fs.fsync
      fn (path: string) -> result[void, string]
      + flushes a file's contents to durable storage
      # filesystem
  std.encoding
    std.encoding.encode_u64
      fn (value: u64) -> bytes
      + encodes an unsigned 64-bit integer in little-endian
      # encoding
    std.encoding.decode_u64
      fn (data: bytes) -> result[u64, string]
      - returns error when data is shorter than 8 bytes
      # encoding
  std.text
    std.text.tokenize
      fn (text: string) -> list[string]
      + splits text into lowercased word tokens dropping punctuation
      # text

graphdb
  graphdb.open
    fn (dir: string) -> result[graphdb_state, string]
    + opens or creates a database rooted at the given directory
    - returns error when the directory cannot be accessed
    # storage
    -> std.fs.read_all
  graphdb.begin_tx
    fn (db: graphdb_state) -> tx_handle
    + starts a new transaction with a private write buffer
    # transactions
  graphdb.add_node
    fn (tx: tx_handle, label: string, props: map[string, string]) -> node_id
    + stages insertion of a node with the given label and properties
    # mutation
  graphdb.add_edge
    fn (tx: tx_handle, from: node_id, to: node_id, kind: string, props: map[string, string]) -> edge_id
    + stages insertion of a directed edge between two nodes
    # mutation
  graphdb.commit
    fn (tx: tx_handle) -> result[void, string]
    + atomically applies the staged mutations and fsyncs the write-ahead log
    - returns error when any staged edge references a missing node
    # transactions
    -> std.fs.append
    -> std.fs.fsync
    -> std.encoding.encode_u64
  graphdb.rollback
    fn (tx: tx_handle) -> void
    + discards all staged mutations for the transaction
    # transactions
  graphdb.index_text
    fn (db: graphdb_state, node: node_id, field: string, text: string) -> void
    + tokenizes text and updates a phrase index mapping tokens to node positions
    # indexing
    -> std.text.tokenize
  graphdb.phrase_search
    fn (db: graphdb_state, phrase: string) -> list[node_id]
    + returns nodes whose indexed text contains the given consecutive token sequence
    - returns empty list when no node matches
    # search
    -> std.text.tokenize
  graphdb.parse_query
    fn (source: string) -> result[query_ast, string]
    + parses a SQL-like query into an abstract syntax tree
    - returns error on syntax failure
    # query_parsing
  graphdb.execute_query
    fn (db: graphdb_state, ast: query_ast) -> result[list[map[string, string]], string]
    + executes a parsed query and returns matched rows
    - returns error when the query references an unknown label
    # query_execution
