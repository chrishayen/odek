# Requirement: "automatic summarization of text documents and HTML pages"

Extractive summarization: strip HTML, tokenize into sentences, score with TF, return the top-N sentences in their original order.

std
  std.html
    std.html.strip_tags
      fn (html: string) -> string
      + returns the concatenated text content of an HTML document
      + collapses adjacent whitespace
      # parsing
  std.text
    std.text.sentence_split
      fn (text: string) -> list[string]
      + splits on sentence-terminating punctuation followed by whitespace
      # text_processing
    std.text.word_tokenize
      fn (text: string) -> list[string]
      + returns lowercase word tokens with punctuation stripped
      # text_processing
    std.text.is_stopword
      fn (token: string) -> bool
      + returns true when the token is in a common stopword list
      # text_processing

summarize
  summarize.from_html
    fn (html: string, max_sentences: i32) -> list[string]
    + strips HTML then delegates to from_text
    # pipeline
    -> std.html.strip_tags
  summarize.from_text
    fn (text: string, max_sentences: i32) -> list[string]
    + returns the top-ranked sentences in source order
    - returns all sentences when the document has fewer than max_sentences
    # pipeline
    -> std.text.sentence_split
  summarize.score_sentences
    fn (sentences: list[string]) -> list[f64]
    + returns a TF-based salience score per sentence excluding stopwords
    # scoring
    -> std.text.word_tokenize
    -> std.text.is_stopword
  summarize.select_top
    fn (sentences: list[string], scores: list[f64], k: i32) -> list[string]
    + returns the k highest-scoring sentences preserving input order
    - returns the empty list when k <= 0
    # selection
