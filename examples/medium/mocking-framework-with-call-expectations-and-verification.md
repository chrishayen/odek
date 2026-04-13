# Requirement: "a mocking framework with call expectations and verification"

Creates mock objects with expected calls, records actual calls, and verifies at the end.

std: (all units exist)

mock
  mock.new_mock
    @ (name: string) -> mock_state
    + creates a mock with the given identifier
    # construction
  mock.expect
    @ (m: mock_state, method: string, args: list[value], returns: list[value]) -> mock_state
    + appends an expected call with argument matchers and the values to return
    # expectation
  mock.expect_times
    @ (m: mock_state, method: string, args: list[value], returns: list[value], times: i32) -> mock_state
    + records that the call must happen exactly times times
    # expectation
  mock.call
    @ (m: mock_state, method: string, args: list[value]) -> result[tuple[list[value], mock_state], string]
    + matches the call against the next unsatisfied expectation and returns its outputs
    - returns error on unexpected method, wrong argument count, or argument mismatch
    # invocation
  mock.verify
    @ (m: mock_state) -> result[void, string]
    + returns ok when every expectation has been satisfied the required number of times
    - returns error listing methods that were under- or over-called
    # verification
  mock.reset
    @ (m: mock_state) -> mock_state
    + clears recorded calls and expectations on the mock
    # lifecycle
