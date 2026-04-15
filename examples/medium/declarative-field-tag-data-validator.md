# Requirement: "a data validation library driven by declarative field tags"

Users describe expected constraints on each field of a record; the library evaluates them and returns a list of violations. No reflection primitive is assumed — the caller supplies field values as a typed map.

std: (all units exist)

validate
  validate.parse_rules
    fn (tag: string) -> result[list[rule], string]
    + returns individual rules parsed from a comma-separated tag string
    - returns error when a rule name is unknown
    # parsing
  validate.required
    fn (name: string, value: field_value) -> optional[violation]
    + returns a violation when the value is empty or absent
    - returns none when the value is non-empty
    # rule
  validate.min_length
    fn (name: string, value: field_value, min: i32) -> optional[violation]
    + returns a violation when the string value is shorter than min
    # rule
  validate.max_length
    fn (name: string, value: field_value, max: i32) -> optional[violation]
    + returns a violation when the string value is longer than max
    # rule
  validate.numeric_range
    fn (name: string, value: field_value, lo: f64, hi: f64) -> optional[violation]
    + returns a violation when a numeric value is outside [lo, hi]
    - returns a violation when the value is not numeric
    # rule
  validate.pattern
    fn (name: string, value: field_value, regex: string) -> optional[violation]
    + returns a violation when a string value does not match the pattern
    # rule
  validate.one_of
    fn (name: string, value: field_value, allowed: list[string]) -> optional[violation]
    + returns a violation when the value is not in the allowed set
    # rule
  validate.check
    fn (schema: map[string, string], record: map[string, field_value]) -> list[violation]
    + returns all violations across all fields in the schema
    + returns an empty list when every field satisfies its rules
    # orchestration
