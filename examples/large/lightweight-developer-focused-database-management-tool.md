# Requirement: "a database management toolkit that manages connections, introspects schemas, and runs queries across relational backends"

The project layer owns connection management, schema introspection, and query execution against a pluggable backend. std carries generic SQL parsing, connection-string parsing, and a result-set iterator primitive.

std
  std.url
    std.url.parse
      @ (raw: string) -> result[url_parts, string]
      + parses scheme, host, port, path, and query
      - returns error on malformed input
      # parsing
  std.sql
    std.sql.tokenize
      @ (source: string) -> list[sql_token]
      + returns the tokens of a SQL string (keywords, identifiers, literals, punctuation)
      # parsing
    std.sql.classify
      @ (source: string) -> string
      + returns the statement class: "select", "insert", "update", "delete", "ddl", or "other"
      # parsing

dbmgr
  dbmgr.new
    @ () -> manager_state
    + creates an empty manager with no registered connections
    # construction
  dbmgr.register_connection
    @ (m: manager_state, id: string, conn_string: string) -> result[void, string]
    + parses the connection string and stores it under id
    - returns error when id is already registered
    - returns error when the connection string is invalid
    # connection
    -> std.url.parse
  dbmgr.remove_connection
    @ (m: manager_state, id: string) -> bool
    + removes the connection and returns true if it existed
    # connection
  dbmgr.list_connections
    @ (m: manager_state) -> list[connection_info]
    + returns all registered connections with id, backend kind, and host
    # connection
  dbmgr.list_schemas
    @ (m: manager_state, id: string) -> result[list[string], string]
    + returns schema names for the identified connection
    - returns error when the connection id is unknown
    # introspection
  dbmgr.list_tables
    @ (m: manager_state, id: string, schema: string) -> result[list[string], string]
    + returns table names in the given schema
    - returns error when schema does not exist
    # introspection
  dbmgr.describe_table
    @ (m: manager_state, id: string, schema: string, table: string) -> result[list[column_info], string]
    + returns column name, type, nullability, and default
    - returns error when table does not exist
    # introspection
  dbmgr.execute
    @ (m: manager_state, id: string, sql: string, params: list[string]) -> result[query_result, string]
    + runs the parameterized SQL and returns rows or a row count
    - returns error when the connection id is unknown
    - returns error when the SQL is rejected by the backend
    # execution
    -> std.sql.classify
    -> std.sql.tokenize
  dbmgr.explain
    @ (m: manager_state, id: string, sql: string) -> result[string, string]
    + returns the backend's query plan as a string
    - returns error when the connection id is unknown
    # execution
