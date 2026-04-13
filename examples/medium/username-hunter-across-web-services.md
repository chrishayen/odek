# Requirement: "a library that checks whether a username exists across a set of web services"

Each target service is described by a URL template and a detection rule. The library fetches each URL and classifies the result.

std
  std.http
    std.http.get
      @ (url: string) -> result[http_response, string]
      + performs a GET and returns status and body
      - returns error on network failure
      # networking
  std.text
    std.text.contains
      @ (haystack: string, needle: string) -> bool
      + returns true when needle occurs in haystack
      # text

username_hunter
  username_hunter.new_target
    @ (name: string, url_template: string, detection: detection_rule) -> target
    + creates a target service descriptor
    ? url_template contains "{username}" which is substituted at check time
    # construction
  username_hunter.check_one
    @ (t: target, username: string) -> result[check_result, string]
    + returns found when the detection rule matches the response
    + returns not_found when the rule indicates absence
    - returns error when the HTTP call fails
    # probing
    -> std.http.get
    -> std.text.contains
  username_hunter.check_all
    @ (targets: list[target], username: string) -> list[check_result]
    + checks each target and collects results; failures become errored entries
    # batch
  username_hunter.found_sites
    @ (results: list[check_result]) -> list[string]
    + returns the names of sites where the username was found
    # filtering
