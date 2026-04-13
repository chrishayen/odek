# Requirement: "a library that exposes shell commands as HTTP endpoints"

Callers register path-to-command mappings; the library turns an incoming HTTP request into a process invocation and formats the result as a response.

std
  std.http
    std.http.parse_request_line
      @ (raw: string) -> result[request_line, string]
      + parses method, path, and query from an HTTP request line
      - returns error when the line does not have three whitespace-separated parts
      # http
    std.http.format_response
      @ (status: i32, body: bytes) -> bytes
      + formats a minimal HTTP/1.1 response with Content-Length
      # http
  std.process
    std.process.spawn_and_wait
      @ (program: string, args: list[string], stdin: bytes) -> result[process_result, string]
      + runs the program to completion and returns stdout, stderr, and exit code
      - returns error when the program cannot be launched
      # process

cmdserve
  cmdserve.new_router
    @ () -> router_state
    + creates an empty command router
    # construction
  cmdserve.register
    @ (router: router_state, path: string, program: string, args: list[string]) -> router_state
    + associates an HTTP path with a command specification
    # registration
  cmdserve.dispatch
    @ (router: router_state, raw_request: string) -> bytes
    + parses the request, runs the matching command, and returns a formatted HTTP response
    - returns a 404 response when no registered path matches
    - returns a 500 response when the command fails to launch
    # dispatch
    -> std.http.parse_request_line
    -> std.process.spawn_and_wait
    -> std.http.format_response
