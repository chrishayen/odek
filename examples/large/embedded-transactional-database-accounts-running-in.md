# Requirement: "an embedded transactional database of accounts with multithreaded access"

The project layer exposes account operations with ACID transactions. std provides primitives for persistent storage, synchronization, and write-ahead logging.

std
  std.fs
    std.fs.append_bytes
      @ (path: string, data: bytes) -> result[void, string]
      + appends data to the file, creating it if absent
      - returns error when the parent directory does not exist
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of the file
      - returns error when the file does not exist
      # filesystem
    std.fs.fsync
      @ (path: string) -> result[void, string]
      + flushes file contents to durable storage
      # filesystem
  std.sync
    std.sync.new_rwlock
      @ () -> rwlock_state
      + creates an unlocked reader-writer lock
      # concurrency
    std.sync.read_lock
      @ (lock: rwlock_state) -> rwlock_state
      + acquires a shared read lock, blocking if a writer holds it
      # concurrency
    std.sync.write_lock
      @ (lock: rwlock_state) -> rwlock_state
      + acquires an exclusive write lock
      # concurrency
    std.sync.unlock
      @ (lock: rwlock_state) -> rwlock_state
      + releases the lock held by the current thread
      # concurrency
  std.encoding
    std.encoding.encode_u64_le
      @ (value: u64) -> bytes
      + encodes a u64 as 8 little-endian bytes
      # encoding
    std.encoding.decode_u64_le
      @ (data: bytes) -> result[u64, string]
      + decodes 8 little-endian bytes into a u64
      - returns error when input is shorter than 8 bytes
      # encoding

accounts_db
  accounts_db.open
    @ (data_path: string, wal_path: string) -> result[db_state, string]
    + opens an existing database or creates a new one at the given paths
    + replays the write-ahead log to recover uncommitted state
    - returns error when paths are unreadable
    # lifecycle
    -> std.fs.read_all
    -> std.sync.new_rwlock
  accounts_db.create_account
    @ (db: db_state, account_id: string, initial_balance: i64) -> result[db_state, string]
    + creates a new account with the given opening balance
    - returns error when the account id already exists
    - returns error when initial balance is negative
    # account_management
  accounts_db.get_balance
    @ (db: db_state, account_id: string) -> result[i64, string]
    + returns the committed balance for the account
    - returns error when the account does not exist
    ? uses a shared read lock so readers do not block each other
    # query
    -> std.sync.read_lock
    -> std.sync.unlock
  accounts_db.begin_tx
    @ (db: db_state) -> tx_state
    + opens a new transaction with a snapshot of current balances
    ? transactions are optimistic; conflicts are detected at commit time
    # transactions
  accounts_db.tx_transfer
    @ (tx: tx_state, from_id: string, to_id: string, amount: i64) -> result[tx_state, string]
    + records a transfer within the transaction without writing it yet
    - returns error when the source balance would go negative
    - returns error when either account does not exist in the snapshot
    # transactions
  accounts_db.commit
    @ (db: db_state, tx: tx_state) -> result[db_state, string]
    + writes the transaction to the WAL, fsyncs, then applies to the main store
    - returns error when a concurrent transaction modified an involved account
    ? holds the writer lock only during WAL append and apply
    # transactions
    -> std.sync.write_lock
    -> std.sync.unlock
    -> std.fs.append_bytes
    -> std.fs.fsync
    -> std.encoding.encode_u64_le
  accounts_db.rollback
    @ (tx: tx_state) -> void
    + discards the transaction snapshot without touching the store
    # transactions
  accounts_db.checkpoint
    @ (db: db_state) -> result[db_state, string]
    + flushes the in-memory store to the data file and truncates the WAL
    - returns error when the fsync fails
    # durability
    -> std.fs.append_bytes
    -> std.fs.fsync
    -> std.encoding.decode_u64_le
  accounts_db.close
    @ (db: db_state) -> result[void, string]
    + performs a final checkpoint and releases resources
    # lifecycle
