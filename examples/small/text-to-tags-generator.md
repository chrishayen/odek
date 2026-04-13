# Requirement: "a library that produces a set of tags from a given text source"

Tokenizes text, strips stop words, and returns the top-k keywords by frequency.

std
  std.text
    std.text.tokenize_words
      @ (text: string) -> list[string]
      + splits on non-letter boundaries and lowercases each token
      # text
    std.text.is_stopword
      @ (word: string) -> bool
      + returns true for common short English stop words
      # text

tagger
  tagger.extract
    @ (text: string, max_tags: i32) -> list[string]
    + returns up to max_tags keywords ordered by frequency then alphabetically
    + ignores tokens shorter than three characters
    - returns an empty list for empty text
    # extraction
    -> std.text.tokenize_words
    -> std.text.is_stopword
  tagger.extract_weighted
    @ (text: string, max_tags: i32) -> list[tuple[string, f64]]
    + returns (tag, normalized_frequency) pairs summing to 1
    # extraction
    -> std.text.tokenize_words
    -> std.text.is_stopword
