# Requirement: "a library for parsing, formatting, storing, and validating international phone numbers"

Entry points cover the full phone-number lifecycle over a region metadata table.

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

intlphone
  intlphone.normalize_digits
    @ (raw: string) -> string
    + removes separators and converts unicode digits to ASCII 0-9
    # normalization
    -> std.strings.trim
    -> std.unicode.is_digit
    -> std.unicode.digit_value
  intlphone.parse
    @ (raw: string, default_region: string) -> result[intl_number, string]
    + recognizes '+' prefix for explicit country codes
    + otherwise interprets as a national number in default_region
    - returns error when no country code can be determined
    # parsing
    -> std.strings.starts_with
  intlphone.is_valid
    @ (num: intl_number) -> bool
    + returns true when length and leading digits match the region rules
    # validation
  intlphone.is_possible
    @ (num: intl_number) -> bool
    + returns true when the length falls in the region's possible range
    # validation
  intlphone.number_type
    @ (num: intl_number) -> string
    + returns "mobile", "fixed_line", "voip", "toll_free", or "unknown"
    # classification
  intlphone.region_for
    @ (num: intl_number) -> string
    + returns the primary region for the number's country code
    # classification
  intlphone.format_e164
    @ (num: intl_number) -> string
    + returns "+<country><national>" with no separators
    # formatting
  intlphone.format_international
    @ (num: intl_number) -> string
    + returns a space-grouped international form
    # formatting
  intlphone.format_national
    @ (num: intl_number) -> string
    + returns the national form with the region's separators
    # formatting
  intlphone.format_rfc3966
    @ (num: intl_number) -> string
    + returns a "tel:+..." URI form
    # formatting
  intlphone.store
    @ (num: intl_number) -> string
    + returns a canonical string suitable for durable storage
    # storage
  intlphone.load
    @ (stored: string) -> result[intl_number, string]
    + parses a canonical stored form back into a number
    - returns error when the stored form is malformed
    # storage
