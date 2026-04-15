# Requirement: "a human name parser that breaks a name into its individual components"

Splits a raw name into title, first, middle, last, suffix, and nickname using classification tables.

std
  std.strings
    std.strings.split
      fn (s: string, sep: string) -> list[string]
      + splits s on every occurrence of sep
      # strings
    std.strings.trim
      fn (s: string) -> string
      + strips leading and trailing ASCII whitespace
      # strings
    std.strings.to_lower
      fn (s: string) -> string
      + returns an ASCII-lowercased copy
      # strings

humanname
  humanname.parse
    fn (raw: string) -> human_name
    + extracts title, first, middle, last, suffix, and nickname
    + handles "Last, First" comma-ordered input
    ? returns empty strings for absent components rather than none
    # parsing
    -> std.strings.split
    -> std.strings.trim
  humanname.is_title
    fn (token: string) -> bool
    + returns true for recognized honorifics like "Dr.", "Mrs."
    # classification
    -> std.strings.to_lower
  humanname.is_suffix
    fn (token: string) -> bool
    + returns true for generational or credential suffixes like "Jr.", "PhD"
    # classification
    -> std.strings.to_lower
  humanname.is_last_name_particle
    fn (token: string) -> bool
    + returns true for particles like "van", "de", "von"
    # classification
    -> std.strings.to_lower
  humanname.extract_nickname
    fn (raw: string) -> tuple[string, string]
    + returns (raw_without_nickname, nickname) when a quoted or parenthesized nickname is present
    # extraction
  humanname.format
    fn (n: human_name) -> string
    + reassembles components into a display form
    # formatting
