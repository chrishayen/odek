# Requirement: "a recursive web content discovery library"

Given a base URL and a wordlist, probe derived URLs and report ones that respond successfully, recursing into discovered directories.

std
  std.http
    std.http.get_status
      @ (url: string) -> result[i32, string]
      + returns the HTTP status code of a GET request
      - returns error on network failure
      # http
  std.strings
    std.strings.trim
      @ (s: string) -> string
      + returns s with leading and trailing whitespace removed
      # strings
    std.strings.join
      @ (parts: list[string], sep: string) -> string
      + joins parts with the given separator
      # strings

discovery
  discovery.new
    @ (base_url: string, wordlist: list[string], max_depth: i32) -> discovery_state
    + returns a scanner rooted at base_url with the given wordlist and recursion limit
    # construction
  discovery.candidate_urls
    @ (state: discovery_state, current: string) -> list[string]
    + returns base urls for each word joined under current
    # enumeration
    -> std.strings.join
  discovery.is_directory_status
    @ (status: i32) -> bool
    + returns true for 200, 301, 302, 401, 403
    - returns false for 404 and 5xx
    # classification
  discovery.probe_url
    @ (url: string) -> result[probe_result, string]
    + returns a probe_result containing url and status
    - returns error when the fetch fails
    # probe
    -> std.http.get_status
  discovery.scan
    @ (state: discovery_state) -> list[probe_result]
    + returns every url that responds with a non-404 status
    + recurses into successful results whose status indicates a directory, up to max_depth
    ? failed fetches are silently skipped
    # orchestration
    -> discovery.candidate_urls
    -> discovery.probe_url
    -> discovery.is_directory_status
