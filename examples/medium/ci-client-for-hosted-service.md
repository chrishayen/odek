# Requirement: "a client for a hosted continuous integration service"

Typed wrappers around build, project, and artifact endpoints of a CI service.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + issues an HTTP request with headers and body and returns status and body
      - returns error on network failure
      # http
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.parse_array_of_objects
      fn (raw: string) -> result[list[map[string, string]], string]
      + parses a JSON array of objects into a list of maps
      - returns error when the root is not an array
      # serialization

ci_client
  ci_client.new
    fn (base_url: string, token: string) -> ci_client_state
    + creates a client that authenticates requests with the given token
    # construction
  ci_client.list_projects
    fn (state: ci_client_state) -> result[list[map[string, string]], string]
    + returns project summaries visible to the token
    - returns error when the token is invalid
    # listing
    -> std.http.request
    -> std.json.parse_array_of_objects
  ci_client.get_build
    fn (state: ci_client_state, project: string, build_num: i64) -> result[map[string, string], string]
    + returns metadata for a single build
    - returns error when build_num does not exist
    # inspection
    -> std.http.request
    -> std.json.parse_object
  ci_client.trigger_build
    fn (state: ci_client_state, project: string, branch: string) -> result[i64, string]
    + triggers a new build on the given branch and returns the new build number
    - returns error when the branch does not exist
    # triggering
    -> std.http.request
    -> std.json.parse_object
  ci_client.list_artifacts
    fn (state: ci_client_state, project: string, build_num: i64) -> result[list[map[string, string]], string]
    + returns artifacts produced by the given build
    # listing
    -> std.http.request
    -> std.json.parse_array_of_objects
