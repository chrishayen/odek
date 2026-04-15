# Requirement: "a general-purpose function mocking library"

A registry that lets test code substitute functions by identifier and record invocations. Call sites look up their replacement through the registry.

std: (all units exist)

mock
  mock.new_registry
    fn () -> registry_state
    + returns an empty registry with no installed mocks and no recorded calls
    # construction
  mock.install
    fn (reg: registry_state, id: string, replacement: any) -> registry_state
    + binds an identifier to a replacement value returned by dispatch
    # registration
  mock.uninstall
    fn (reg: registry_state, id: string) -> registry_state
    + removes the replacement for an identifier; dispatch falls through afterwards
    # registration
  mock.dispatch
    fn (reg: registry_state, id: string, args: list[any]) -> tuple[registry_state, optional[any]]
    + records the call under the identifier and returns the installed replacement when present
    - returns none for the result when no replacement is installed
    # dispatch
  mock.calls
    fn (reg: registry_state, id: string) -> list[list[any]]
    + returns the recorded argument lists for every dispatch targeting the identifier
    # introspection
  mock.reset
    fn (reg: registry_state) -> registry_state
    + clears all installed replacements and recorded calls
    # lifecycle
