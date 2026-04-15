# Requirement: "a collection of tools and datasets for natural language processing of a logographic language"

Tokenization, stopword filtering, and dataset lookup for a writing system without word boundaries.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when path does not exist
      # filesystem
  std.string
    std.string.split
      fn (s: string, sep: string) -> list[string]
      + splits on the given separator
      # strings

nlp
  nlp.segment
    fn (text: string, dictionary: list[string]) -> list[string]
    + greedy longest-match segmentation into tokens
    + handles text with no matching words by returning single codepoints
    # tokenization
  nlp.load_dictionary
    fn (path: string) -> result[list[string], string]
    + loads a newline-delimited word list from disk
    - returns error when the file cannot be read
    # dataset_loading
    -> std.fs.read_all
    -> std.string.split
  nlp.load_stopwords
    fn (path: string) -> result[list[string], string]
    + loads a stopword set
    # dataset_loading
    -> std.fs.read_all
    -> std.string.split
  nlp.filter_stopwords
    fn (tokens: list[string], stopwords: list[string]) -> list[string]
    + returns tokens that are not in the stopword set
    # filtering
  nlp.count_terms
    fn (tokens: list[string]) -> map[string, i64]
    + tallies token frequencies
    # frequency
