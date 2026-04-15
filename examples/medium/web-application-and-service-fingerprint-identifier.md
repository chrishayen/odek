# Requirement: "a library for identifying web application and service fingerprints from response data"

Matches HTTP response headers, bodies, and status codes against a rule set to identify the software serving the response.

std
  std.regex
    std.regex.match
      fn (pattern: string, text: string) -> bool
      + returns true when the pattern matches anywhere in the text
      # pattern_matching
  std.hash
    std.hash.md5_hex
      fn (data: bytes) -> string
      + returns the hex-encoded md5 digest
      # hashing

fingerprint
  fingerprint.load_rules
    fn (raw: string) -> result[list[fingerprint_rule], string]
    + parses a rule set mapping patterns to product names
    - returns error on malformed rule entries
    # rule_loading
  fingerprint.match_headers
    fn (headers: map[string, string], rule: fingerprint_rule) -> bool
    + returns true when any header pattern in the rule matches
    # matching
    -> std.regex.match
  fingerprint.match_body
    fn (body: string, rule: fingerprint_rule) -> bool
    + returns true when the rule's body pattern matches
    # matching
    -> std.regex.match
  fingerprint.match_favicon
    fn (favicon: bytes, rule: fingerprint_rule) -> bool
    + returns true when the favicon hash equals the rule's expected hash
    # matching
    -> std.hash.md5_hex
  fingerprint.identify
    fn (response: http_response, rules: list[fingerprint_rule]) -> list[string]
    + returns product names for every rule that matches the response
    # identification
