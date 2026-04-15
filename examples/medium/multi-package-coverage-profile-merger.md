# Requirement: "a multi-package coverage profile merger"

Scans a project tree for per-package coverage profile files, merges them into a single combined profile, and summarizes line coverage.

std
  std.fs
    std.fs.walk
      fn (root: string) -> list[string]
      + returns all file paths under root, recursively
      # filesystem
    std.fs.read_text
      fn (path: string) -> result[string, string]
      + reads a file as UTF-8 text
      - returns error when the file does not exist
      # filesystem
    std.fs.write_text
      fn (path: string, content: string) -> result[void, string]
      + writes text to a file, creating or truncating it
      # filesystem

coverage
  coverage.parse_profile
    fn (raw: string) -> result[list[coverage_block], string]
    + parses a coverage profile into block records
    - returns error on malformed lines
    ? each block has file, start_line, end_line, num_statements, count
    # parsing
  coverage.discover_profiles
    fn (root: string, filename_suffix: string) -> list[string]
    + returns paths under root whose name ends with the suffix
    # discovery
    -> std.fs.walk
  coverage.load_profiles
    fn (paths: list[string]) -> result[list[coverage_block], string]
    + concatenates parsed blocks from all profiles
    - returns error on the first unreadable or malformed file
    # loading
    -> std.fs.read_text
    -> coverage.parse_profile
  coverage.merge_blocks
    fn (blocks: list[coverage_block]) -> list[coverage_block]
    + sums counts for blocks with identical file and line range
    ? block order follows first appearance
    # merging
  coverage.write_profile
    fn (blocks: list[coverage_block], path: string) -> result[void, string]
    + writes a merged profile in the standard text format
    # output
    -> std.fs.write_text
  coverage.summarize
    fn (blocks: list[coverage_block]) -> coverage_summary
    + returns total statements, covered statements, and percent
    # summary
