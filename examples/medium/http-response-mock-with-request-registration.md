# Requirement: "a library for mocking HTTP responses by registering expected requests"

Registers rules matching method, URL, and optional body, returns canned responses, and records which rules were hit.

std: (all units exist)

httpmock
  httpmock.new
    fn () -> mock_state
    + creates an empty mock registry
    # construction
  httpmock.register
    fn (state: mock_state, method: string, url_pattern: string, status: i32, body: string) -> mock_state
    + registers a rule that responds with status and body when method and url match
    ? url_pattern matches exact strings; "*" is a wildcard
    # registration
  httpmock.register_body_match
    fn (state: mock_state, method: string, url_pattern: string, body_substring: string, status: i32, body: string) -> mock_state
    + registers a rule that also requires the request body to contain body_substring
    # registration
  httpmock.serve
    fn (state: mock_state, method: string, url: string, body: string) -> result[mock_response, string]
    + returns the matched canned response and increments its hit counter
    - returns error when no rule matches
    # serving
  httpmock.hits_for
    fn (state: mock_state, method: string, url_pattern: string) -> i32
    + returns the number of times the rule was matched
    # verification
  httpmock.assert_all_called
    fn (state: mock_state) -> result[void, list[string]]
    + returns ok when every registered rule was matched at least once
    - returns the list of uncalled rule descriptions otherwise
    # verification
