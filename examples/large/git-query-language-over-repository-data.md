# Requirement: "a SQL-like query language over git repository data"

A library that parses a small SQL dialect and executes SELECT queries against commits, branches, and tags read from a repository on disk.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the path does not exist or is unreadable
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns entry names in the directory
      - returns error when the path is not a directory
      # filesystem
  std.hash
    std.hash.sha1_hex
      @ (data: bytes) -> string
      + returns the lowercase hex sha1 of data
      # hashing
  std.compress
    std.compress.zlib_inflate
      @ (data: bytes) -> result[bytes, string]
      + decompresses a zlib stream
      - returns error on corrupt input
      # compression

git_query
  git_query.open_repo
    @ (path: string) -> result[repo_state, string]
    + opens a repository rooted at path and indexes refs
    - returns error when .git is not found under path
    # repository
    -> std.fs.read_all
    -> std.fs.list_dir
  git_query.load_commit
    @ (r: repo_state, oid: string) -> result[commit_record, string]
    + returns a commit's author, committer, message, parents and tree oid
    - returns error when the oid does not resolve to a commit object
    # object_loading
    -> std.compress.zlib_inflate
    -> std.hash.sha1_hex
  git_query.walk_commits
    @ (r: repo_state, start_ref: string, limit: i32) -> result[list[commit_record], string]
    + returns commits reachable from start_ref in topological order, up to limit
    - returns error when start_ref cannot be resolved
    # history
  git_query.tokenize
    @ (source: string) -> result[list[token], string]
    + splits a query string into keywords, identifiers, literals and punctuation
    - returns error on an unterminated string literal
    # lexing
  git_query.parse_select
    @ (tokens: list[token]) -> result[select_ast, string]
    + parses SELECT columns FROM table WHERE expr ORDER BY expr LIMIT n
    - returns error when FROM is missing
    - returns error on an unknown keyword in the WHERE expression
    # parsing
  git_query.bind_table
    @ (r: repo_state, ast: select_ast) -> result[row_source, string]
    + binds FROM commits|branches|tags to a row iterator over the repo
    - returns error when the table name is not recognized
    # binding
  git_query.evaluate_where
    @ (src: row_source, ast: select_ast) -> result[row_source, string]
    + filters rows using the parsed WHERE expression
    - returns error when the expression references an unknown column
    # evaluation
  git_query.project_and_order
    @ (src: row_source, ast: select_ast) -> result[list[row], string]
    + selects the requested columns, applies ORDER BY and LIMIT
    - returns error when an ORDER BY column is not in the projection
    # projection
  git_query.execute
    @ (r: repo_state, source: string) -> result[list[row], string]
    + parses and runs a query against the opened repository
    - returns error on any parse, bind, or evaluation failure
    # query_execution
    -> git_query.tokenize
    -> git_query.parse_select
    -> git_query.bind_table
    -> git_query.evaluate_where
    -> git_query.project_and_order
