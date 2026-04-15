# Requirement: "a client library for a continuous integration server"

Thin client for talking to a remote CI server over HTTP. Transport and credentials are injected so tests can stub them.

std
  std.http
    std.http.get
      fn (url: string, headers: map[string, string]) -> result[http_response, string]
      + issues a GET and returns status, headers, and body
      - returns error on network failure
      # http
    std.http.post
      fn (url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + issues a POST with the given body
      - returns error on network failure
      # http
  std.encoding
    std.encoding.base64_encode
      fn (data: bytes) -> string
      + encodes bytes as standard base64 with padding
      # encoding

ci_client
  ci_client.new
    fn (base_url: string, user: string, token: string) -> ci_client_state
    + constructs a client with base URL and basic auth credentials
    # construction
    -> std.encoding.base64_encode
  ci_client.list_jobs
    fn (client: ci_client_state) -> result[list[job_summary], string]
    + returns names and last-build statuses of top-level jobs
    - returns error on non-2xx response
    # jobs
    -> std.http.get
  ci_client.trigger_build
    fn (client: ci_client_state, job: string, params: map[string, string]) -> result[i64, string]
    + queues a build and returns the queue item id
    - returns error when the job does not exist
    # builds
    -> std.http.post
  ci_client.get_build_status
    fn (client: ci_client_state, job: string, build_number: i64) -> result[build_status, string]
    + returns the current state of a specific build
    - returns error when the build does not exist
    # builds
    -> std.http.get
  ci_client.fetch_build_log
    fn (client: ci_client_state, job: string, build_number: i64) -> result[string, string]
    + returns the console log as a string
    # builds
    -> std.http.get
