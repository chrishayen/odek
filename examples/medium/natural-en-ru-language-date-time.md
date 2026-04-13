# Requirement: "a natural-language date and time parser with pluggable rules for multiple locales"

Rules match substrings and produce relative or absolute offsets against a reference instant. Locales contribute their own rule sets.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns the current unix time in seconds
      # time
    std.time.add_seconds
      @ (epoch_seconds: i64, delta: i64) -> i64
      + returns epoch_seconds + delta
      # time

when
  when.parser_new
    @ () -> parser_state
    + creates a parser with no registered locales
    # construction
  when.register_locale
    @ (state: parser_state, locale: string, rules: list[rule]) -> parser_state
    + registers the given rule set under the locale code
    ? a rule is an opaque token-matching callable with a tag and priority
    # configuration
  when.register_rule
    @ (state: parser_state, locale: string, rule: rule) -> result[parser_state, string]
    + appends a single rule to an existing locale
    - returns error when the locale is not registered
    # configuration
  when.tokenize
    @ (text: string) -> list[token]
    + splits text into lowercase word tokens with byte offsets
    # lexing
  when.match_rules
    @ (state: parser_state, locale: string, tokens: list[token]) -> list[match]
    + returns every non-overlapping match sorted by start offset then priority
    # matching
  when.resolve
    @ (match: match, reference_epoch: i64) -> i64
    + computes the absolute unix time for a match using reference_epoch
    # resolution
    -> std.time.add_seconds
  when.parse
    @ (state: parser_state, text: string, locale: string, reference_epoch: i64) -> result[list[resolved], string]
    + returns every recognized date/time expression with its absolute epoch and character span
    - returns error when the locale is not registered
    # parsing
    -> std.time.now_seconds
