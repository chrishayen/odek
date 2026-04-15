# Requirement: "a library that fetches trending open-source projects from a code hosting API"

Given a language filter and a time window, queries a hosting provider's search endpoint and returns a ranked list.

std
  std.http
    std.http.get
      fn (url: string, headers: map[string, string]) -> result[string, string]
      + performs an HTTP GET and returns the response body
      - returns error on non-2xx status
      # http_client
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization

trending
  trending.build_query
    fn (language: string, since_days: i32) -> string
    + returns a provider-ready search query for repos in a language created within since_days
    # query_construction
  trending.fetch
    fn (base_url: string, query: string) -> result[list[map[string, string]], string]
    + returns a list of repository records sorted by star count descending
    - returns error when the provider rejects the request
    # fetch
    -> std.http.get
    -> std.json.parse_object
