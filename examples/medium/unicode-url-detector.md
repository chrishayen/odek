# Requirement: "a unicode-aware URL detector"

Scans text for URLs, emails, and bare domains, returning ranges. Unicode handling goes through std.

std
  std.unicode
    std.unicode.decode_utf8
      @ (data: bytes) -> result[list[i32], string]
      + decodes bytes into codepoints
      - returns error on invalid utf-8
      # unicode
    std.unicode.is_letter
      @ (cp: i32) -> bool
      + returns true when the codepoint is in any Letter category
      # unicode
    std.unicode.is_digit
      @ (cp: i32) -> bool
      + returns true when the codepoint is a decimal digit
      # unicode

linkify
  linkify.find_links
    @ (text: string) -> list[link_match]
    + returns every detected link with its byte offset, length, and kind
    + recognizes http, https, ftp schemes and bare domains with recognized TLDs
    - returns empty list when no link is found
    # scanning
    -> std.unicode.decode_utf8
  linkify.is_scheme_char
    @ (cp: i32) -> bool
    + returns true for letters, digits, '+', '-', '.' as allowed in a URI scheme
    # classification
    -> std.unicode.is_letter
    -> std.unicode.is_digit
  linkify.match_domain
    @ (codepoints: list[i32], start: i32) -> optional[i32]
    + returns the end offset of a valid domain starting at start, or none
    ? a domain is one or more labels separated by '.', ending with a known TLD
    # domain_parsing
    -> std.unicode.is_letter
  linkify.match_email
    @ (codepoints: list[i32], start: i32) -> optional[i32]
    + returns the end offset of a valid email starting at start, or none
    # email_parsing
    -> std.unicode.is_letter
  linkify.register_tld
    @ (state: linkify_state, tld: string) -> linkify_state
    + adds a TLD to the recognized set
    # configuration
