# Requirement: "a human name parser that splits a name into its individual parts"

Splits a raw name string into title, first, middle, last, suffix, and nickname using rule tables.

std
  std.strings
    std.strings.split
      @ (s: string, sep: string) -> list[string]
      + splits s on every occurrence of sep
      # strings
    std.strings.to_lower
      @ (s: string) -> string
      + returns an ASCII-lowercased copy
      # strings
    std.strings.trim
      @ (s: string) -> string
      + strips leading and trailing ASCII whitespace
      # strings

nameparse
  nameparse.parse
    @ (raw: string) -> parsed_name
    + extracts title, first, middle, last, suffix, and nickname
    + handles "Last, First" ordering
    ? returns empty strings for parts that are absent rather than none
    # parsing
    -> std.strings.split
    -> std.strings.trim
  nameparse.is_title
    @ (token: string) -> bool
    + returns true for recognized honorifics like "Dr.", "Mrs.", "Sir"
    # classification
    -> std.strings.to_lower
  nameparse.is_suffix
    @ (token: string) -> bool
    + returns true for generational or credential suffixes like "Jr.", "III", "PhD"
    # classification
    -> std.strings.to_lower
  nameparse.is_last_name_prefix
    @ (token: string) -> bool
    + returns true for particles like "van", "de", "del", "von"
    # classification
    -> std.strings.to_lower
  nameparse.extract_nickname
    @ (raw: string) -> tuple[string, string]
    + returns (raw_without_nickname, nickname) when the raw string contains a quoted or parenthesized nickname
    # extraction
  nameparse.format_full
    @ (n: parsed_name) -> string
    + reassembles the parts into a canonical display form
    # formatting
