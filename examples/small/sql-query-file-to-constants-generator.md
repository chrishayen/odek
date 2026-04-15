# Requirement: "a tool that generates constants from files containing raw sql queries"

Reads a directory of query files and emits a source file that exposes each query as a named string constant.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when the file cannot be opened
      # filesystem
    std.fs.list_dir
      fn (dir: string) -> result[list[string], string]
      + returns the names of entries directly inside the directory
      - returns error when the directory does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes the contents to the given path, replacing any existing file
      # filesystem

sql_consts
  sql_consts.collect_queries
    fn (dir: string) -> result[list[query_entry], string]
    + returns one entry per file in the directory, each holding the file stem as a name and the trimmed file contents as the query
    - returns error when any query file cannot be read
    # loading
    -> std.fs.list_dir
    -> std.fs.read_all
  sql_consts.render_constants
    fn (entries: list[query_entry]) -> string
    + returns source text declaring one string constant per entry, named after the file stem
    ? constant names are uppercased versions of the file stem
    # codegen
  sql_consts.generate_file
    fn (input_dir: string, output_path: string) -> result[void, string]
    + collects queries and writes the rendered source file to disk
    - returns error when collection or writing fails
    # orchestration
    -> std.fs.write_all
