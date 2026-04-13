# Requirement: "a library for recording and replaying http interactions in tests"

On the first run, calls to the http client are recorded to a cassette file. On subsequent runs, matching requests are replayed from the cassette. Project layer owns the cassette and matching; std provides file I/O, json, and http.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads the entire file at path as a utf-8 string
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: string) -> result[void, string]
      + writes data to path, creating or truncating the file
      # filesystem
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses any JSON value
      - returns error on malformed input
      # serialization
    std.json.encode_value
      @ (value: json_value) -> string
      + encodes a JSON value as a string
      # serialization
  std.http
    std.http.send
      @ (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an http request and returns the response
      - returns error on network failure
      # networking

vcr
  vcr.load_cassette
    @ (path: string) -> result[cassette_state, string]
    + reads a cassette file from disk into memory
    + returns an empty cassette when the file does not exist
    # persistence
    -> std.fs.read_all
    -> std.json.parse_value
  vcr.save_cassette
    @ (c: cassette_state, path: string) -> result[void, string]
    + serializes the cassette and writes it to path
    # persistence
    -> std.json.encode_value
    -> std.fs.write_all
  vcr.match_request
    @ (c: cassette_state, method: string, url: string, body: bytes) -> optional[http_response]
    + returns the recorded response whose (method, url, body-hash) matches
    - returns none when nothing matches
    # matching
  vcr.record_interaction
    @ (c: cassette_state, method: string, url: string, body: bytes, resp: http_response) -> cassette_state
    + appends a request/response pair to the cassette
    # recording
  vcr.request
    @ (c: cassette_state, method: string, url: string, headers: map[string, string], body: bytes) -> result[tuple[http_response, cassette_state], string]
    + returns a matched recording when present; otherwise performs the real request and records it
    - returns error when there is no match and the network call fails
    # playback
    -> std.http.send
