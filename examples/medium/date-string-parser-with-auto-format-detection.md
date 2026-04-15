# Requirement: "parse date strings without knowing the format in advance"

A date-string parser that probes a set of candidate formats and returns the first unambiguous match as a unix timestamp.

std: (all units exist)

dateparse
  dateparse.parse
    fn (raw: string) -> result[i64, string]
    + returns a unix seconds timestamp on the first matching format
    - returns error when no candidate format matches
    - returns error when the input is empty
    # parsing
  dateparse.parse_in_zone
    fn (raw: string, zone_offset_seconds: i32) -> result[i64, string]
    + returns a unix timestamp interpreting naive inputs in the given zone offset
    - returns error when no candidate format matches
    # parsing
  dateparse.detect_format
    fn (raw: string) -> result[string, string]
    + returns the canonical format descriptor that matched the input
    - returns error when the input is ambiguous between two distinct formats
    - returns error when the input matches no known format
    # introspection
  dateparse.try_formats
    fn (raw: string, candidates: list[string]) -> result[tuple[i64, string], string]
    + returns (timestamp, matching_format) for the first candidate that parses raw
    - returns error when none of the candidates parse
    # parsing
  dateparse.normalize_whitespace
    fn (raw: string) -> string
    + returns the input with leading/trailing whitespace stripped and internal runs collapsed to single spaces
    # preprocessing
