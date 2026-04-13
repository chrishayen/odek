# Requirement: "an http client mocking library"

Registers canned responses for (method, url) pairs and serves them to intercepted calls.

std: (all units exist)

http_mock
  http_mock.new
    @ () -> mock_registry
    + returns an empty registry with no registered responses
    # construction
  http_mock.register
    @ (registry: mock_registry, method: string, url: string, status: i32, body: string) -> mock_registry
    + stores a canned response keyed by (method, url)
    + later registrations overwrite earlier ones for the same key
    # registration
  http_mock.match
    @ (registry: mock_registry, method: string, url: string) -> result[tuple[i32, string], string]
    + returns (status, body) for a registered (method, url) pair
    - returns error when no matching response has been registered
    # dispatch
