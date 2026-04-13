# Requirement: "a library that loads test fixtures from SQL files before a test and cleans them up afterwards"

Setup executes the fixture statements; teardown reverses them by truncating the touched tables.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns file contents as text
      - returns error when the path does not exist
      # filesystem
  std.db
    std.db.exec
      @ (db: db_handle, sql: string) -> result[void, string]
      + executes a statement that returns no rows
      - returns error on syntax or constraint failure
      # database

testsql
  testsql.split_statements
    @ (sql: string) -> list[string]
    + returns each top-level statement, stripping trailing semicolons
    + ignores semicolons inside quoted strings
    # parsing
  testsql.extract_target_tables
    @ (statements: list[string]) -> list[string]
    + returns the table names touched by INSERT or UPDATE statements
    ? used later to decide what to truncate during cleanup
    # analysis
  testsql.setup
    @ (db: db_handle, fixture_path: string) -> result[list[string], string]
    + reads the file, splits it, runs each statement, returns the touched tables
    - returns error on the first failing statement
    # setup
    -> std.fs.read_all
    -> std.db.exec
  testsql.teardown
    @ (db: db_handle, tables: list[string]) -> result[void, string]
    + truncates each table in reverse order
    - returns error when a truncate fails
    # teardown
