# Requirement: "a parser for human-readable dates in many languages"

Tokenizes free-form date strings, matches them against language-specific patterns, and returns a normalized timestamp.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
    std.time.compose
      @ (year: i32, month: i32, day: i32, hour: i32, minute: i32, second: i32) -> result[i64, string]
      + returns unix seconds for the given UTC date/time components
      - returns error when any component is out of range
      # time

dateparser
  dateparser.new
    @ () -> dateparser_state
    + returns a parser with the default set of languages enabled
    # construction
  dateparser.enable_language
    @ (state: dateparser_state, lang: string) -> dateparser_state
    + enables pattern tables for the given language code
    # registration
  dateparser.tokenize
    @ (state: dateparser_state, input: string) -> list[token]
    + splits input into normalized words, numbers, and punctuation
    # tokenization
  dateparser.match_relative
    @ (state: dateparser_state, tokens: list[token]) -> optional[i64]
    + recognizes phrases like "yesterday" or "3 days ago" and returns the timestamp
    - returns none when no relative phrase matches
    # parsing
    -> std.time.now_seconds
  dateparser.match_absolute
    @ (state: dateparser_state, tokens: list[token]) -> optional[i64]
    + recognizes explicit dates such as "March 5, 2024" in any enabled language
    - returns none when no absolute pattern matches
    # parsing
    -> std.time.compose
  dateparser.parse
    @ (state: dateparser_state, input: string) -> result[i64, string]
    + returns a unix timestamp for the given human-readable date
    - returns error when the input cannot be interpreted
    # entry_point
