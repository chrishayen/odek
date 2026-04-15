# Requirement: "a multilingual text segmentation library"

Tokenizes text into words using a dictionary with frequency weights, supporting languages with and without whitespace word boundaries.

std
  std.unicode
    std.unicode.is_whitespace
      fn (cp: i32) -> bool
      + returns true for unicode whitespace code points
      # unicode
    std.unicode.is_cjk
      fn (cp: i32) -> bool
      + returns true for CJK unified ideograph ranges
      # unicode
    std.unicode.decode_utf8
      fn (data: bytes) -> result[list[i32], string]
      + decodes a UTF-8 byte sequence into code points
      - returns error on invalid UTF-8
      # unicode

segmenter
  segmenter.new
    fn () -> segmenter_state
    + returns a segmenter with an empty dictionary
    # construction
  segmenter.load_dictionary
    fn (state: segmenter_state, entries: list[tuple[string, f64]]) -> segmenter_state
    + adds (term, frequency) entries to the dictionary
    ? frequency is a weight used during maximum-likelihood segmentation
    # dictionary
  segmenter.cut_whitespace
    fn (state: segmenter_state, text: string) -> list[string]
    + splits on whitespace runs for languages with explicit word boundaries
    # segmentation
    -> std.unicode.is_whitespace
    -> std.unicode.decode_utf8
  segmenter.cut_dag
    fn (state: segmenter_state, text: string) -> list[string]
    + builds a directed acyclic graph of dictionary matches and returns the most-probable segmentation
    ? uses dynamic programming over log-frequencies
    # segmentation
    -> std.unicode.decode_utf8
  segmenter.cut
    fn (state: segmenter_state, text: string) -> list[string]
    + dispatches to whitespace or DAG segmentation based on script detection
    + returns the concatenation when no dictionary matches are found
    # segmentation
    -> std.unicode.is_cjk
    -> std.unicode.decode_utf8
