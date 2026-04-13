# Requirement: "a simple way to generate mocks for multiple purposes"

Builds canned-response mocks that record invocations. The mock is a state object; call records are stored as tuples of (method name, args).

std: (all units exist)

mockgen
  mockgen.new
    @ () -> mock_state
    + creates an empty mock with no canned responses and no recorded calls
    # construction
  mockgen.when
    @ (state: mock_state, method: string, args: list[string], returns: string) -> mock_state
    + registers that calling `method` with the given args returns the given value
    # setup
  mockgen.call
    @ (state: mock_state, method: string, args: list[string]) -> tuple[result[string, string], mock_state]
    + looks up a canned response and records the invocation
    - returns error when no canned response matches
    # invocation
  mockgen.calls_of
    @ (state: mock_state, method: string) -> list[list[string]]
    + returns every recorded arg list for the given method in call order
    # inspection
