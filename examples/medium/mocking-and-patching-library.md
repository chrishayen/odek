# Requirement: "a mocking and patching library"

Mocks record calls and return canned values; patches temporarily replace a named attribute in a registry.

std: (all units exist)

mocking
  mocking.new_mock
    fn (return_value: optional[string]) -> mock_state
    + creates a fresh mock with zero recorded calls
    + a none return_value means the mock returns none when called
    # construction
  mocking.call
    fn (state: mock_state, args: list[string]) -> tuple[optional[string], mock_state]
    + records the argument list and returns the canned value
    + call count increments by one per invocation
    # invocation
  mocking.call_count
    fn (state: mock_state) -> i32
    + returns the number of times call was invoked
    # introspection
  mocking.called_with
    fn (state: mock_state, args: list[string]) -> bool
    + returns true when any recorded call matches the given argument list
    - returns false when no recorded call matches
    # assertion
  mocking.reset
    fn (state: mock_state) -> mock_state
    + clears recorded calls while preserving the return value
    # lifecycle
  mocking.new_registry
    fn () -> registry_state
    + creates an empty attribute registry
    # construction
  mocking.set
    fn (reg: registry_state, name: string, value: string) -> registry_state
    + stores a value under the given name
    # registry
  mocking.patch
    fn (reg: registry_state, name: string, replacement: string) -> tuple[registry_state, string]
    + replaces the value at name with replacement and returns the original
    - returns an empty original when name was not previously set
    # patching
  mocking.unpatch
    fn (reg: registry_state, name: string, original: string) -> registry_state
    + restores the original value at name
    # patching
