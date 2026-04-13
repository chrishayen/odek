# Requirement: "a natural language processing toolkit"

Core NLP primitives: tokenization, sentence splitting, part-of-speech tagging, stemming, stopword filtering, and frequency analysis.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads a file into a string
      - returns error when the file does not exist
      # filesystem
  std.text
    std.text.to_lower
      @ (s: string) -> string
      + returns the lowercase form of a string
      # text
    std.text.split_whitespace
      @ (s: string) -> list[string]
      + splits on any run of whitespace, dropping empties
      # text

nlptk
  nlptk.tokenize
    @ (text: string) -> list[string]
    + splits text into word-level tokens, stripping trailing punctuation
    + keeps contractions as single tokens
    # tokenization
    -> std.text.split_whitespace
  nlptk.split_sentences
    @ (text: string) -> list[string]
    + splits text into sentences using terminal punctuation
    + treats "Mr." and similar abbreviations as non-terminal
    # segmentation
  nlptk.pos_tag
    @ (tokens: list[string]) -> list[tuple[string,string]]
    + returns each token paired with a coarse part-of-speech tag
    ? tagger uses a dictionary lookup with unknown-word fallback
    # tagging
  nlptk.stem
    @ (word: string) -> string
    + returns a suffix-stripped stem for an English word
    + is idempotent for already-stemmed words
    # morphology
  nlptk.load_stopwords
    @ (path: string) -> result[list[string], string]
    + reads a stopword list, one word per line
    - returns error when the file does not exist
    # filtering
    -> std.fs.read_all
  nlptk.remove_stopwords
    @ (tokens: list[string], stopwords: list[string]) -> list[string]
    + returns the tokens not present in the stopword list
    # filtering
    -> std.text.to_lower
  nlptk.frequency
    @ (tokens: list[string]) -> map[string, i32]
    + returns a token-to-count map
    # statistics
  nlptk.top_n
    @ (freq: map[string,i32], n: i32) -> list[tuple[string,i32]]
    + returns the n most frequent entries sorted descending
    # statistics
  nlptk.ngrams
    @ (tokens: list[string], n: i32) -> list[list[string]]
    + returns every contiguous window of n tokens
    - returns an empty list when n exceeds the token count
    # ngrams
