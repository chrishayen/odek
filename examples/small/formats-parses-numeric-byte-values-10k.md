# Requirement: "format and parse human-readable byte sizes like 10K, 2M, 3G"

Bidirectional conversion between byte counts and short human-readable strings.

std: (all units exist)

bytesize
  bytesize.format
    @ (bytes: i64) -> string
    + returns a compact form like "10K", "2M", "3G" using 1024-based units
    + returns "0B" for 0
    + returns "1023B" for 1023
    + returns "1K" for 1024
    ? chooses the largest unit that yields an integer or short decimal
    # formatting
  bytesize.format_precision
    @ (bytes: i64, decimals: i32) -> string
    + like format but keeps the given number of decimal places
    + returns "1.50K" for 1536 with decimals=2
    # formatting
  bytesize.parse
    @ (s: string) -> result[i64, string]
    + parses "10K", "2M", "3G", "512B", "1.5G"
    + accepts both uppercase and lowercase suffixes
    - returns error on unknown suffix
    - returns error on non-numeric prefix
    # parsing
