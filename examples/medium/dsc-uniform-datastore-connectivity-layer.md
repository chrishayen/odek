# Requirement: "a datastore connectivity layer exposing a uniform interface over relational, document, and file-backed stores"

A single open/query/exec surface that dispatches to the registered driver for the store's scheme.

std: (all units exist)

dsc
  dsc.register_driver
    fn (scheme: string, driver: driver) -> void
    + registers a driver for a URL scheme
    # registry
  dsc.open
    fn (url: string) -> result[connection, string]
    + parses the scheme from the URL and delegates to the matching driver
    - returns error when no driver is registered for the scheme
    # connect
  dsc.query
    fn (conn: connection, statement: string, params: list[value]) -> result[list[map[string, value]], string]
    + returns rows as maps from column name to value
    - returns error when the underlying driver rejects the statement
    # query
  dsc.exec
    fn (conn: connection, statement: string, params: list[value]) -> result[i64, string]
    + returns the number of affected rows
    - returns error when the underlying driver rejects the statement
    # command
  dsc.close
    fn (conn: connection) -> result[void, string]
    + releases resources held by the connection
    # teardown
