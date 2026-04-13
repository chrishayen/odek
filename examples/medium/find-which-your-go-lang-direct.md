# Requirement: "a library to find which of your direct source-hosted dependencies are susceptible to chainjacking"

Inspects a list of direct dependencies hosted on a source-control platform and flags any whose upstream owner account has been deleted or renamed, leaving the namespace reclaimable.

std
  std.http
    std.http.get
      @ (url: string) -> result[http_response, string]
      + performs a GET and returns status, headers, and body
      - returns error on network failure
      # networking
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a json_value
      - returns error on malformed JSON
      # serialization
    std.json.get_field
      @ (value: json_value, field: string) -> optional[json_value]
      + returns the named field of a JSON object
      - returns none when the field is absent
      # serialization

chainjack
  chainjack.parse_dependency_path
    @ (path: string) -> result[repo_ref, string]
    + parses "host.example/owner/repo" into host, owner, and repo
    - returns error when path has fewer than three segments
    # parsing
  chainjack.list_direct_dependencies
    @ (manifest: string) -> result[list[repo_ref], string]
    + returns every direct dependency listed in a module manifest
    - returns error on malformed manifest
    # discovery
  chainjack.check_owner_exists
    @ (host: string, owner: string) -> result[bool, string]
    + returns true when the owner account resolves on the host
    - returns false when the owner returns 404
    - returns error on non-404 failures
    # probing
    -> std.http.get
  chainjack.check_repo_exists
    @ (ref: repo_ref) -> result[bool, string]
    + returns true when the repository resolves on the host
    # probing
    -> std.http.get
    -> std.json.parse
  chainjack.classify
    @ (ref: repo_ref, owner_exists: bool, repo_exists: bool) -> chainjack_status
    + returns vulnerable when owner is missing (namespace reclaimable)
    + returns suspicious when owner exists but repo is missing
    + returns safe when both resolve
    # classification
  chainjack.audit
    @ (manifest: string) -> result[list[chainjack_finding], string]
    + returns a finding per direct dependency with status and details
    - returns error when the manifest cannot be parsed
    # orchestration
    -> std.json.get_field
  chainjack.format_report
    @ (findings: list[chainjack_finding]) -> string
    + renders a human-readable report grouped by status
    # rendering
