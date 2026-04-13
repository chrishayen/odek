# Requirement: "a trigram-indexed code search engine"

Source files are broken into 3-character grams; a posting list maps each gram to the documents containing it. Queries intersect posting lists for the grams in the query string and then verify matches line by line.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents as bytes
      - returns error when the file is missing
      # filesystem
    std.fs.walk
      @ (root: string) -> result[list[string], string]
      + returns every file path under root recursively
      # filesystem
  std.text
    std.text.split_lines
      @ (data: string) -> list[string]
      + splits input at LF and CRLF boundaries
      # text

code_search
  code_search.extract_trigrams
    @ (content: string) -> list[u32]
    + returns the set of 3-byte grams packed into u32
    ? strings shorter than 3 bytes contribute no grams
    # tokenization
  code_search.new_index
    @ () -> index_state
    + returns an empty index with no documents
    # construction
  code_search.index_file
    @ (state: index_state, path: string, content: string) -> index_state
    + adds the document and inserts it into the posting list of each of its trigrams
    # indexing
  code_search.index_tree
    @ (state: index_state, root: string) -> result[index_state, string]
    + walks root and indexes every readable file
    - returns error on filesystem failure
    # bulk_indexing
    -> std.fs.walk
    -> std.fs.read_all
  code_search.candidate_docs
    @ (state: index_state, query: string) -> list[string]
    + returns documents whose posting lists contain all grams of the query
    + returns every document when the query is shorter than 3 bytes
    # retrieval
  code_search.verify_matches
    @ (path: string, content: string, query: string) -> list[match]
    + returns (line_number, line_text) pairs containing the query substring
    # verification
    -> std.text.split_lines
  code_search.search
    @ (state: index_state, query: string) -> result[list[file_matches], string]
    + returns matches across the index, intersecting posting lists before verifying
    - returns error when a candidate file cannot be re-read
    # query
    -> std.fs.read_all
