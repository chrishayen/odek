# Requirement: "a copy-paste detector for source code"

Tokenize source, compute rolling fingerprints over k-grams, and report duplicate spans across files.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads an entire file as UTF-8 text
      - returns error when the file does not exist or cannot be read
      # filesystem
  std.hash
    std.hash.fnv64
      @ (data: bytes) -> u64
      + returns a 64-bit FNV-1a hash
      # hashing

cpd
  cpd.tokenize
    @ (source: string) -> list[token]
    + splits source into identifier, literal, operator, and whitespace tokens
    + strips comments and collapses whitespace runs
    # lexing
  cpd.fingerprint
    @ (tokens: list[token], k: i32) -> list[fingerprint]
    + emits a rolling fingerprint for each k-gram of significant tokens
    + each fingerprint carries its start and end token index
    # hashing
    -> std.hash.fnv64
  cpd.index_file
    @ (idx: cpd_index, file: string, source: string, k: i32) -> cpd_index
    + adds a file's fingerprints to an in-memory index keyed by hash
    # indexing
  cpd.find_duplicates
    @ (idx: cpd_index, min_tokens: i32) -> list[clone]
    + returns clone groups where the matched span is at least min_tokens long
    + extends matching k-grams into the longest contiguous duplicate span
    # detection
  cpd.scan_paths
    @ (paths: list[string], k: i32, min_tokens: i32) -> result[list[clone], string]
    + convenience entry point: reads files, indexes them, and returns clone groups
    # pipeline
    -> std.fs.read_all
