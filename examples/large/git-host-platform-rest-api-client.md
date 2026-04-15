# Requirement: "a client library for interacting with a source code hosting platform's rest api"

Library layer — request building and response parsing — that a CLI or other app could wrap. No cli runes.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs the request and returns status, headers, body
      - returns error on network failure
      # http
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.parse_array
      fn (raw: string) -> result[list[string], string]
      + parses a JSON array into a list of raw elements
      - returns error on non-array root
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.url
    std.url.build
      fn (base: string, path: string, query: map[string, string]) -> string
      + joins base + path and appends a url-encoded query string
      # url

git_host_client
  git_host_client.new
    fn (base_url: string, token: string) -> client_state
    + creates a client that sends the token as a bearer header on every call
    # construction
  git_host_client.list_projects
    fn (c: client_state, page: i32, per_page: i32) -> result[list[string], string]
    + returns json for each project on the given page
    - returns error when status is not 2xx
    # projects
    -> std.url.build
    -> std.http.request
    -> std.json.parse_array
  git_host_client.get_project
    fn (c: client_state, project_id: string) -> result[map[string, string], string]
    + returns project fields keyed by name
    - returns error when project does not exist
    # projects
    -> std.http.request
    -> std.json.parse_object
  git_host_client.list_issues
    fn (c: client_state, project_id: string, state: string) -> result[list[string], string]
    + returns issues filtered by state ("open" or "closed")
    - returns error when state is not a known value
    # issues
    -> std.url.build
    -> std.http.request
    -> std.json.parse_array
  git_host_client.create_issue
    fn (c: client_state, project_id: string, title: string, body: string) -> result[string, string]
    + returns the id of the created issue
    - returns error when title is empty
    # issues
    -> std.json.encode_object
    -> std.http.request
    -> std.json.parse_object
  git_host_client.list_merge_requests
    fn (c: client_state, project_id: string) -> result[list[string], string]
    + returns open merge requests for the project
    # merge_requests
    -> std.http.request
    -> std.json.parse_array
  git_host_client.merge_request_status
    fn (c: client_state, project_id: string, mr_id: string) -> result[string, string]
    + returns "open", "merged", or "closed"
    - returns error when the mr id is unknown
    # merge_requests
    -> std.http.request
    -> std.json.parse_object
