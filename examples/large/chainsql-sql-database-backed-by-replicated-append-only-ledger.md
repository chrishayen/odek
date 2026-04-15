# Requirement: "a SQL database backed by a replicated append-only ledger"

Clients submit SQL statements; each statement is serialized into a ledger block, replicated across peers via consensus, and applied to a local SQL engine.

std
  std.crypto
    std.crypto.sha256
      fn (data: bytes) -> bytes
      + returns the SHA-256 digest of data (32 bytes)
      # cryptography
    std.crypto.ed25519_sign
      fn (private_key: bytes, message: bytes) -> bytes
      + returns an Ed25519 signature of message
      # cryptography
    std.crypto.ed25519_verify
      fn (public_key: bytes, message: bytes, signature: bytes) -> bool
      + returns true when the signature is valid
      - returns false on any error or mismatch
      # cryptography
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.sql
    std.sql.open_in_memory
      fn () -> sql_engine
      + opens a fresh in-memory SQL engine
      # storage
    std.sql.execute
      fn (e: sql_engine, statement: string) -> result[i64, string]
      + executes a DDL or DML statement and returns the affected row count
      - returns error on syntax or constraint violation
      # storage
    std.sql.query
      fn (e: sql_engine, statement: string) -> result[list[map[string, string]], string]
      + executes a SELECT and returns the rows as string-keyed maps
      # storage

chainsql
  chainsql.new_node
    fn (node_id: string, private_key: bytes) -> node
    + creates a node with an empty ledger and a fresh SQL engine
    # construction
    -> std.sql.open_in_memory
  chainsql.sign_statement
    fn (n: node, sql: string) -> signed_statement
    + produces a signed statement carrying node id, timestamp, and signature
    # authoring
    -> std.crypto.ed25519_sign
    -> std.time.now_millis
  chainsql.verify_statement
    fn (s: signed_statement, public_key: bytes) -> bool
    + returns true when the signature matches the author key
    # validation
    -> std.crypto.ed25519_verify
  chainsql.propose_block
    fn (n: node, statements: list[signed_statement]) -> block
    + bundles statements into a block with parent hash and merkle root
    # block_production
    -> std.crypto.sha256
  chainsql.merkle_root
    fn (statements: list[signed_statement]) -> bytes
    + returns the merkle root of statement hashes
    # hashing
    -> std.crypto.sha256
  chainsql.append_block
    fn (n: node, b: block) -> result[node, string]
    + validates the block and appends it to the local ledger
    - returns error when the parent hash does not match the current tip
    - returns error when any statement signature is invalid
    # ledger
  chainsql.apply_block
    fn (n: node, b: block) -> result[node, string]
    + executes each statement against the local SQL engine in order
    + halts and returns error on the first statement that fails
    # state_machine
    -> std.sql.execute
  chainsql.submit
    fn (n: node, sql: string) -> result[node, string]
    + signs, proposes, appends, and applies a single statement
    # entry_point
  chainsql.query
    fn (n: node, sql: string) -> result[list[map[string, string]], string]
    + runs a read-only query against the local engine without touching the ledger
    # read_path
    -> std.sql.query
  chainsql.ledger_height
    fn (n: node) -> i64
    + returns the current block height
    # introspection
  chainsql.ledger_tip_hash
    fn (n: node) -> string
    + returns the hex-encoded hash of the most recent block
    # introspection
    -> std.encoding.hex_encode
