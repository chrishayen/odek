# Requirement: "a library for running SQL against production databases with access controls, audit logging, and shareable query links"

A gated SQL gateway: statements are parsed, checked against per-user ACLs, executed through a pluggable database host, and logged. Results can be persisted and retrieved by share token.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.random
    std.random.token
      @ (n_bytes: i32) -> string
      + returns a url-safe random token of approximately the requested byte length
      # random
  std.sql
    std.sql.parse_statement
      @ (raw: string) -> result[parsed_statement, string]
      + parses a SQL statement and classifies it (select, insert, update, delete, ddl) with referenced tables
      - returns error on malformed SQL
      # parsing
  std.json
    std.json.encode_rows
      @ (columns: list[string], rows: list[list[string]]) -> string
      + encodes a tabular result set as JSON
      # serialization

sqlgate
  sqlgate.new_gateway
    @ (db: database_host, audit: audit_sink) -> gateway_state
    + creates a gateway bound to a database host and an audit sink
    # construction
  sqlgate.define_role
    @ (state: gateway_state, role: string) -> gateway_state
    + registers a role with no permissions
    # acl
  sqlgate.grant
    @ (state: gateway_state, role: string, table: string, ops: list[string]) -> gateway_state
    + grants the given operations on a table to a role
    # acl
  sqlgate.assign_role
    @ (state: gateway_state, user: string, role: string) -> gateway_state
    + associates a user with a role
    # acl
  sqlgate.check_access
    @ (state: gateway_state, user: string, stmt: parsed_statement) -> result[void, string]
    + allows the request when the user's role grants the statement's operation on every referenced table
    - returns error naming the first missing permission
    # authorization
  sqlgate.execute
    @ (state: gateway_state, user: string, raw_sql: string) -> result[query_result, string]
    + parses the statement, checks access, runs it via the database host, and emits an audit record
    - returns error when parsing or authorization fails, without touching the database
    # execution
    -> std.sql.parse_statement
    -> std.time.now_millis
  sqlgate.create_share
    @ (state: gateway_state, result: query_result, ttl_seconds: i64) -> share_record
    + persists the result and mints a random share token with expiry
    # sharing
    -> std.random.token
    -> std.time.now_millis
  sqlgate.resolve_share
    @ (state: gateway_state, token: string) -> result[query_result, string]
    + returns the stored result for an unexpired token
    - returns error when the token is unknown or expired
    # sharing
    -> std.time.now_millis
  sqlgate.serialize_result
    @ (result: query_result) -> string
    + encodes a result set as JSON for transport or storage
    # serialization
    -> std.json.encode_rows
  sqlgate.audit_entries
    @ (state: gateway_state, since_ms: i64) -> list[audit_entry]
    + returns audit records newer than the given timestamp
    # audit
