# Requirement: "a code coverage library that aggregates data from subprocesses"

A library that instruments a coverage map, merges data files written by subprocesses, and computes coverage summaries.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns entries in a directory
      - returns error when the directory does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating or truncating it
      # filesystem
  std.encoding
    std.encoding.json_decode_any
      fn (raw: string) -> result[json_value, string]
      + parses arbitrary JSON into a dynamic value
      - returns error on invalid JSON
      # serialization
    std.encoding.json_encode_any
      fn (value: json_value) -> string
      + encodes a dynamic JSON value
      # serialization

coverage
  coverage.new_map
    fn () -> coverage_map
    + creates an empty coverage map
    # construction
  coverage.record_statement
    fn (cmap: coverage_map, file: string, line: i32) -> coverage_map
    + increments the execution count for a statement location
    # recording
  coverage.record_branch
    fn (cmap: coverage_map, file: string, line: i32, branch: i32, taken: bool) -> coverage_map
    + records whether a branch at a location was taken
    # recording
  coverage.merge
    fn (a: coverage_map, b: coverage_map) -> coverage_map
    + combines two coverage maps by summing counts and unioning branch sets
    # merging
  coverage.load_from_file
    fn (path: string) -> result[coverage_map, string]
    + reads a serialized coverage map from disk
    - returns error when the file is malformed
    # persistence
    -> std.fs.read_all
    -> std.encoding.json_decode_any
  coverage.save_to_file
    fn (cmap: coverage_map, path: string) -> result[void, string]
    + writes a coverage map to disk as JSON
    # persistence
    -> std.encoding.json_encode_any
    -> std.fs.write_all
  coverage.load_from_directory
    fn (dir: string) -> result[coverage_map, string]
    + merges every coverage file in a directory (one per subprocess)
    - returns error when the directory cannot be read
    # aggregation
    -> std.fs.list_dir
  coverage.summarize
    fn (cmap: coverage_map) -> coverage_summary
    + returns total and covered statement and branch counts
    # reporting
