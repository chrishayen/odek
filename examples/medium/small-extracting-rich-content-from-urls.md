# Requirement: "an oembed-style rich content extractor"

Given a URL, look up the matching provider pattern and fetch the rich-content metadata. HTTP is a std primitive so the extractor can be tested with a fake transport.

std
  std.http
    std.http.get
      @ (url: string) -> result[string, string]
      + returns the response body for 2xx status codes
      - returns error on network failure or non-2xx status
      # http
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
  std.strings
    std.strings.matches_glob
      @ (s: string, pattern: string) -> bool
      + returns true when the string matches the glob pattern
      # strings

rich_content
  rich_content.new_registry
    @ () -> provider_registry
    + returns an empty provider registry
    # construction
  rich_content.register_provider
    @ (reg: provider_registry, url_pattern: string, endpoint: string) -> provider_registry
    + returns a registry with the provider added
    # registration
  rich_content.find_provider
    @ (reg: provider_registry, url: string) -> optional[string]
    + returns the endpoint template whose pattern matches the URL
    - returns none when no pattern matches
    # lookup
    -> std.strings.matches_glob
  rich_content.extract
    @ (reg: provider_registry, url: string) -> result[map[string, string], string]
    + returns the parsed rich-content fields for a URL
    - returns error when no provider matches
    - returns error when the provider request fails
    # extraction
    -> std.http.get
    -> std.json.parse_object
