# Requirement: "an HTTP mocking library for intercepting outgoing requests"

A fluent registry of matchers. Each registered mock specifies what to match (method, url pattern, header constraints) and the response to return.

std: (all units exist)

http_mock
  http_mock.new
    fn () -> mock_registry
    + returns a registry with no registered mocks
    # construction
  http_mock.expect
    fn (registry: mock_registry, method: string, url_pattern: string) -> mock_builder
    + returns a builder scoped to a new mock for (method, url_pattern)
    ? url_pattern supports glob-style wildcards (*)
    # matching
  http_mock.with_header
    fn (builder: mock_builder, name: string, value: string) -> mock_builder
    + adds a required header constraint to the mock under construction
    # matching
  http_mock.with_body
    fn (builder: mock_builder, substring: string) -> mock_builder
    + adds a required body-substring constraint to the mock under construction
    # matching
  http_mock.reply
    fn (builder: mock_builder, status: i32, body: string) -> mock_registry
    + finalizes the mock with the given response and returns the updated registry
    # registration
  http_mock.dispatch
    fn (registry: mock_registry, method: string, url: string, headers: map[string, string], body: string) -> result[tuple[i32, string], string]
    + returns (status, body) from the first mock whose constraints all match the request
    - returns error when no mock matches
    # dispatch
  http_mock.unmatched_count
    fn (registry: mock_registry) -> i32
    + returns the number of registered mocks that have not yet been dispatched
    ? useful for asserting all registered mocks were exercised
    # introspection
