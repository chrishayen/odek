# Requirement: "a file-based HTTP mock server with recording ability"

Routes HTTP requests by matching them against files in a directory and can record incoming requests for replay.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file's raw bytes
      - returns error when the file does not exist
      # io
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file
      # io
    std.fs.list_dir
      @ (dir: string) -> result[list[string], string]
      + returns entry names inside dir
      # io
  std.http
    std.http.parse_request
      @ (raw: bytes) -> result[http_request, string]
      + parses an HTTP/1.1 request
      - returns error on malformed request line
      # parsing
    std.http.serialize_response
      @ (resp: http_response) -> bytes
      + serializes status, headers, and body to bytes
      # serialization

smoke
  smoke.load_mocks
    @ (dir: string) -> result[list[mock], string]
    + reads every file under dir and parses it into a mock entry
    + a mock file encodes method, path, status, headers, and body
    - returns error when a mock file is malformed
    # configuration
    -> std.fs.list_dir
    -> std.fs.read_all
  smoke.match_mock
    @ (mocks: list[mock], method: string, path: string) -> optional[mock]
    + returns the first mock whose method and path match
    + supports wildcard segments like "/users/{id}"
    # matching
  smoke.handle
    @ (mocks: list[mock], request: bytes) -> result[bytes, string]
    + parses the request, finds a matching mock, and serializes its response
    - returns a 404 response when no mock matches
    # routing
    -> std.http.parse_request
    -> std.http.serialize_response
  smoke.record
    @ (dir: string, method: string, path: string, response: http_response) -> result[void, string]
    + writes a new mock file capturing the response
    # recording
    -> std.fs.write_all
