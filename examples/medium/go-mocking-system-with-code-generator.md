# Requirement: "a mocking system that records and replays calls to interface methods"

Runtime recorders that let tests assert on the sequence of method calls a collaborator received, plus expectation matching.

std: (all units exist)

mock
  mock.new
    @ () -> mock_state
    + creates an empty recorder
    # construction
  mock.record_call
    @ (state: mock_state, method: string, args: list[string]) -> mock_state
    + appends a call entry with its method name and stringified arguments
    # recording
  mock.expect
    @ (state: mock_state, method: string, args: list[string], return_value: string) -> mock_state
    + registers a canned return for a future call matching method and args
    # expectation
  mock.dispatch
    @ (state: mock_state, method: string, args: list[string]) -> result[string, string]
    + returns the canned value and records the call when an expectation matches
    - returns error when no expectation matches
    # dispatch
  mock.calls_for
    @ (state: mock_state, method: string) -> list[list[string]]
    + returns argument lists of every recorded call to method in order
    # inspection
  mock.verify
    @ (state: mock_state) -> result[bool, string]
    + returns true when every registered expectation was matched at least once
    - returns error listing the unmet expectations
    # verification
