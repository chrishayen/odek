# Requirement: "a database driver abstraction that exposes a uniform interface over multiple backends"

A registry of drivers and a uniform API that dispatches queries to whichever driver is bound. Individual backends are plugged in by the caller.

std: (all units exist)

db
  db.new_registry
    @ () -> driver_registry
    + creates an empty driver registry
    # construction
  db.register_driver
    @ (registry: driver_registry, name: string, driver: driver_handle) -> driver_registry
    + registers a driver under a scheme name
    - returns registry unchanged when name is already registered
    # registration
  db.open
    @ (registry: driver_registry, dsn: string) -> result[connection_handle, string]
    + parses dsn "scheme://..." and opens a connection using the matching driver
    - returns error when no driver matches the scheme
    # connection
  db.close
    @ (conn: connection_handle) -> result[void, string]
    + releases the connection
    # connection
  db.exec
    @ (conn: connection_handle, sql: string, params: list[db_value]) -> result[exec_summary, string]
    + runs a statement that does not return rows
    - returns error on backend-reported failure
    # statement
  db.query
    @ (conn: connection_handle, sql: string, params: list[db_value]) -> result[row_cursor, string]
    + runs a statement that returns rows
    # query
  db.next_row
    @ (cursor: row_cursor) -> optional[list[db_value]]
    + returns the next row or none when exhausted
    # query
  db.begin_tx
    @ (conn: connection_handle) -> result[tx_handle, string]
    + opens a transaction on conn
    # transaction
  db.commit_tx
    @ (tx: tx_handle) -> result[void, string]
    + commits the transaction
    - returns error when the underlying driver refuses
    # transaction
  db.rollback_tx
    @ (tx: tx_handle) -> result[void, string]
    + rolls the transaction back
    # transaction
