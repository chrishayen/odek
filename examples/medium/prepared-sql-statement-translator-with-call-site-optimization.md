# Requirement: "a library that transparently uses prepared SQL statements without code changes at call sites"

Wraps a database handle. The first call with a given query string prepares and caches a statement; subsequent calls reuse it.

std
  std.hash
    std.hash.string_fnv64
      @ (s: string) -> u64
      + returns a stable 64-bit FNV hash of the string
      # hashing
  std.sync
    std.sync.mutex_new
      @ () -> mutex
      + creates an unlocked mutex
      # concurrency
    std.sync.mutex_with
      @ (m: mutex, body: fn() -> void) -> void
      + runs body while holding the mutex
      # concurrency

prep
  prep.wrap
    @ (db_handle: db_handle) -> prep_state
    + wraps a database handle with an empty statement cache
    # construction
  prep.get_or_prepare
    @ (state: prep_state, query: string) -> result[stmt_handle, string]
    + returns a cached prepared statement for the query, preparing it on first use
    - returns error when the driver fails to prepare
    # caching
    -> std.hash.string_fnv64
    -> std.sync.mutex_with
  prep.exec
    @ (state: prep_state, query: string, args: list[string]) -> result[i64, string]
    + executes the prepared statement with the given bound arguments and returns affected rows
    - returns error when binding or execution fails
    # execution
  prep.query
    @ (state: prep_state, query: string, args: list[string]) -> result[list[map[string, string]], string]
    + runs the prepared statement and returns rows as string-to-string maps
    - returns error when the query fails
    # querying
  prep.close
    @ (state: prep_state) -> result[void, string]
    + finalizes all cached statements
    # teardown
