# Requirement: "a library for validating structured input against declarative rules"

Rules are composable: built-ins cover common checks (required, length, range, regex) and users may plug in custom functions. Validation returns every failure, not just the first.

std
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex, string]
      + compiles a regular expression
      - returns error on malformed pattern
      # regex
    std.regex.is_match
      fn (re: regex, input: string) -> bool
      + returns true when the input matches the regex
      # regex

validator
  validator.required
    fn () -> rule
    + returns a rule that fails when the field value is empty
    # rule
  validator.min_length
    fn (n: i32) -> rule
    + returns a rule that fails when the value has fewer than n characters
    # rule
  validator.max_length
    fn (n: i32) -> rule
    + returns a rule that fails when the value has more than n characters
    # rule
  validator.in_range
    fn (low: f64, high: f64) -> rule
    + returns a rule that fails when a numeric value lies outside [low, high]
    - fails when the value cannot be parsed as a number
    # rule
  validator.matches
    fn (pattern: string) -> result[rule, string]
    + returns a rule that fails when the value does not match the pattern
    - returns error on invalid pattern
    # rule
    -> std.regex.compile
    -> std.regex.is_match
  validator.custom
    fn (fn: fn(string) -> optional[string]) -> rule
    + returns a rule whose failure message is supplied by the caller's function
    # rule
  validator.check_field
    fn (name: string, value: string, rules: list[rule]) -> list[validation_error]
    + returns every rule that fails against the value, tagged with the field name
    + returns empty list when all rules pass
    # validation
  validator.check_record
    fn (record: map[string, string], schema: map[string, list[rule]]) -> list[validation_error]
    + runs the rule list for each field in the schema and aggregates all failures
    + fields present in the record but not in the schema are ignored
    # validation
