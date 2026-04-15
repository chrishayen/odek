# Requirement: "a bridge library that lets a spreadsheet call into a scripting runtime and vice versa"

Exposes a narrow RPC-style surface between a spreadsheet cell environment and a runtime.

std: (all units exist)

bridge
  bridge.new
    fn () -> bridge_state
    + creates a bridge with no registered handlers
    # construction
  bridge.register
    fn (state: bridge_state, name: string, handler: callback) -> bridge_state
    + registers a named handler that accepts a cell-range argument and returns a value
    # registration
  bridge.call_from_sheet
    fn (state: bridge_state, name: string, args: list[string]) -> result[string, string]
    + invokes the named handler and returns its stringified result
    - returns error when the name is not registered
    # sheet_to_runtime
  bridge.write_range
    fn (state: bridge_state, sheet: string, cell: string, values: list[list[string]]) -> result[void, string]
    + pushes a 2D block of values to a named range
    - returns error on malformed cell reference
    # runtime_to_sheet
