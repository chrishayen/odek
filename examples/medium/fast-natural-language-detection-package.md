# Requirement: "a natural language detector"

Scores input text against per-language trigram profiles and returns the best match with a confidence.

std
  std.unicode
    std.unicode.to_lower
      @ (s: string) -> string
      + returns the lowercase form, handling non-ascii codepoints
      # unicode
    std.unicode.normalize_nfc
      @ (s: string) -> string
      + returns the nfc-normalized form of s
      # unicode

langdetect
  langdetect.clean_text
    @ (text: string) -> string
    + lowercases, normalizes, and strips punctuation and digits
    # preprocessing
    -> std.unicode.to_lower
    -> std.unicode.normalize_nfc
  langdetect.trigrams
    @ (text: string) -> map[string, i32]
    + returns a frequency map of character trigrams with sentence padding
    # feature_extraction
  langdetect.load_profiles
    @ (raw: string) -> result[language_profiles, string]
    + parses an embedded per-language trigram frequency table
    - returns error on malformed data
    # model
  langdetect.score
    @ (profiles: language_profiles, sample: map[string, i32]) -> list[language_score]
    + returns languages ranked by cosine similarity against sample
    # scoring
  langdetect.detect
    @ (profiles: language_profiles, text: string) -> optional[language_score]
    + returns the best-scoring language along with confidence
    - returns none when text is too short to score
    # detection
