# Requirement: "a fake backend library that serves responses from configured routes backed by files or literal objects"

Match incoming requests against a route table and return canned responses.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads an entire file as bytes
      - returns error on missing file
      # filesystem

fakeback
  fakeback.new
    @ () -> fakeback_state
    + creates an empty route table
    # construction
  fakeback.add_route
    @ (state: fakeback_state, method: string, path: string, response: fake_response) -> fakeback_state
    + registers a response keyed by method and path
    # registration
  fakeback.respond_from_file
    @ (path: string, status: i32) -> result[fake_response, string]
    + loads response body from disk and pairs it with a status code
    - returns error when the file cannot be read
    # response_building
    -> std.fs.read_all
  fakeback.respond_from_literal
    @ (body: string, status: i32) -> fake_response
    + wraps a literal body and status into a response
    # response_building
  fakeback.handle
    @ (state: fakeback_state, method: string, path: string) -> result[fake_response, string]
    + returns the configured response for the matched route
    - returns error when no route matches
    # dispatch
