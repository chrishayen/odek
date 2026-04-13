# Requirement: "a lightweight streaming SQL engine for real-time data processing"

Parses a SELECT statement, plans it against a declared stream schema, and evaluates it over a sliding window of rows.

std: (all units exist)

stream_sql
  stream_sql.tokenize
    @ (sql: string) -> result[list[token], string]
    + splits the query into keyword, identifier, literal, and punctuation tokens
    - returns error on unterminated string literals
    # lexing
  stream_sql.parse_select
    @ (tokens: list[token]) -> result[select_stmt, string]
    + parses SELECT ... FROM stream [WHERE ...] [GROUP BY ...] [WINDOW size]
    - returns error when required clauses are missing or out of order
    # parsing
  stream_sql.parse_expr
    @ (tokens: list[token], start: i32) -> result[tuple[expr, i32], string]
    + parses comparison, boolean, and arithmetic expressions with precedence
    - returns error on unexpected tokens
    # parsing
  stream_sql.register_stream
    @ (state: engine_state, name: string, columns: list[column_def]) -> result[engine_state, string]
    + declares a stream schema under the given name
    - returns error when name is already registered
    # schema
  stream_sql.plan
    @ (state: engine_state, stmt: select_stmt) -> result[plan, string]
    + resolves column references and returns a plan of projection, filter, group, window ops
    - returns error on unknown columns or stream name
    # planning
  stream_sql.push_row
    @ (state: engine_state, stream: string, row: row) -> result[engine_state, string]
    + appends a row to the named stream's window buffer
    - returns error when stream is not registered
    # ingestion
  stream_sql.evaluate_expr
    @ (e: expr, row: row) -> result[value, string]
    + computes the expression value against the row's columns
    - returns error on type mismatches
    # execution
  stream_sql.execute
    @ (state: engine_state, p: plan) -> result[list[row], string]
    + runs the plan over the current window and returns result rows
    - returns error when the window is empty and no default is configured
    # execution
  stream_sql.aggregate
    @ (rows: list[row], group_keys: list[string], agg: aggregate_spec) -> result[list[row], string]
    + groups rows and computes count, sum, avg, min, max
    - returns error on unknown aggregate name
    # execution
  stream_sql.slide_window
    @ (state: engine_state, stream: string, size: i32) -> engine_state
    + trims the stream's buffer to the most recent size rows
    # windowing
