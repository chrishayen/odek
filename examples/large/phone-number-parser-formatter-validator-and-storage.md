# Requirement: "a phone number parsing, formatting, storage, and validation library"

The project exposes parse/format/validate entry points over a pluggable metadata table keyed by region. std provides generic string and unicode primitives.

std
  std.strings
    std.strings.trim
      @ (s: string) -> string
      + strips leading and trailing ASCII whitespace
      # strings
    std.strings.starts_with
      @ (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings
  std.unicode
    std.unicode.is_digit
      @ (cp: i32) -> bool
      + returns true for unicode decimal digits
      # unicode
    std.unicode.digit_value
      @ (cp: i32) -> optional[i32]
      + returns 0..9 for a decimal-digit code point
      # unicode

phonenum
  phonenum.normalize
    @ (raw: string) -> string
    + strips spaces, dashes, parentheses, and dots
    + converts unicode digits to ASCII 0-9
    # normalization
    -> std.strings.trim
    -> std.unicode.is_digit
    -> std.unicode.digit_value
  phonenum.parse
    @ (raw: string, default_region: string) -> result[phone_number, string]
    + recognizes a leading '+' followed by a country calling code
    + falls back to default_region when input lacks a country prefix
    - returns error when the resulting national number is empty
    - returns error when the country code is unknown
    # parsing
    -> std.strings.starts_with
  phonenum.country_code_for_region
    @ (region: string) -> optional[i32]
    + looks up ISO 3166-1 alpha-2 region to calling code
    # metadata
  phonenum.region_for_country_code
    @ (code: i32) -> optional[string]
    + returns the primary region for a calling code
    # metadata
  phonenum.is_valid
    @ (num: phone_number) -> bool
    + returns true when the national length matches the region's rules
    - returns false when the number is shorter or longer than allowed
    # validation
  phonenum.is_possible
    @ (num: phone_number) -> bool
    + returns true when length falls within the region's possible range
    # validation
  phonenum.format_e164
    @ (num: phone_number) -> string
    + returns "+<country><national>" with no separators
    # formatting
  phonenum.format_international
    @ (num: phone_number) -> string
    + returns a space-separated international form
    # formatting
  phonenum.format_national
    @ (num: phone_number) -> string
    + returns the national form using the region's grouping rules
    # formatting
  phonenum.number_type
    @ (num: phone_number) -> string
    + classifies as "mobile", "fixed_line", "toll_free", or "unknown"
    # classification
  phonenum.store_serialize
    @ (num: phone_number) -> string
    + returns a stable canonical E.164 form for storage
    # storage
  phonenum.store_deserialize
    @ (stored: string) -> result[phone_number, string]
    + parses a canonical stored form back into a phone_number
    - returns error when the stored form is malformed
    # storage
