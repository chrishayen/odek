# Requirement: "a recursive text search tool that also extracts text from rich document formats"

Core search engine plus a registry of extractors that convert structured files to plain text before scanning.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the complete contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns every file path under root, depth-first
      - returns error when root is not a directory
      # filesystem
  std.regex
    std.regex.compile
      @ (pattern: string) -> result[regex, string]
      + returns a compiled pattern
      - returns error on invalid regex syntax
      # regex
    std.regex.find_all
      @ (re: regex, text: string) -> list[match_span]
      + returns all non-overlapping matches as (start, end, line_number) triples
      # regex
  std.archive
    std.archive.list_entries
      @ (data: bytes, format: u8) -> result[list[archive_entry], string]
      + returns entries (0=zip, 1=tar, 2=tar.gz) with name and offset
      - returns error when the archive is malformed
      # archive
    std.archive.read_entry
      @ (data: bytes, format: u8, entry: archive_entry) -> result[bytes, string]
      + returns the decompressed bytes for a single entry
      # archive

multisearch
  multisearch.new
    @ (pattern: string) -> result[searcher_state, string]
    + compiles the pattern and returns a searcher with no registered extractors
    - returns error when pattern is invalid
    # construction
    -> std.regex.compile
  multisearch.register_extractor
    @ (state: searcher_state, extension: string, kind: u8) -> searcher_state
    + associates a file extension with an extractor (0=pdf, 1=epub, 2=office, 3=archive, 4=plain)
    # extractor_registry
  multisearch.extract_pdf
    @ (data: bytes) -> result[string, string]
    + returns the concatenated text of every page
    - returns error on encrypted PDFs
    # extraction
  multisearch.extract_epub
    @ (data: bytes) -> result[string, string]
    + returns the concatenated text of spine items
    - returns error on malformed containers
    # extraction
    -> std.archive.list_entries
    -> std.archive.read_entry
  multisearch.extract_office
    @ (data: bytes) -> result[string, string]
    + returns the concatenated text from the document body parts
    - returns error when the document.xml part is missing
    # extraction
    -> std.archive.list_entries
    -> std.archive.read_entry
  multisearch.extract_archive
    @ (data: bytes, format: u8) -> result[list[tuple[string, string]], string]
    + returns (entry_name, text) pairs for each text-like entry inside the archive
    # extraction
    -> std.archive.list_entries
    -> std.archive.read_entry
  multisearch.search_text
    @ (state: searcher_state, text: string) -> list[match_span]
    + returns matches within a plain-text body
    # search
    -> std.regex.find_all
  multisearch.search_file
    @ (state: searcher_state, path: string) -> result[list[match_span], string]
    + reads the file, runs the extractor for its extension, and returns matches
    - returns error when the file cannot be read
    # search
    -> std.fs.read_all
  multisearch.search_tree
    @ (state: searcher_state, root: string) -> result[list[tuple[string, match_span]], string]
    + walks root and returns (path, match) pairs for every hit in supported files
    - returns error when root is not a directory
    # search
    -> std.fs.walk
