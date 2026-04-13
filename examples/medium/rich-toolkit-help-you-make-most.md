# Requirement: "an interactive computing toolkit for exploring code and data"

A kernel-style evaluation loop backed by a cell store and a simple display protocol. The runtime host supplies an evaluator callback; this library manages state.

std
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time

notebook
  notebook.new
    @ () -> notebook_state
    + creates an empty notebook with counter zero and no cells
    # construction
  notebook.add_cell
    @ (state: notebook_state, kind: u8, source: string) -> tuple[string, notebook_state]
    + appends a cell (kind: 0=code, 1=markdown) and returns its id
    # cell_management
  notebook.update_cell
    @ (state: notebook_state, id: string, source: string) -> notebook_state
    + replaces the source of an existing cell
    - returns unchanged state when id is unknown
    # cell_management
  notebook.delete_cell
    @ (state: notebook_state, id: string) -> notebook_state
    + removes the cell
    # cell_management
  notebook.record_execution
    @ (state: notebook_state, id: string, output: string, error: optional[string]) -> notebook_state
    + stores the result of executing a code cell, increments the execution counter, and records wall-clock time
    # execution
    -> std.time.now_millis
  notebook.cell_display
    @ (state: notebook_state, id: string) -> result[string, string]
    + returns the JSON display payload for a cell (source, output, execution_count)
    - returns error when id is unknown
    # display
    -> std.json.encode_object
  notebook.export_document
    @ (state: notebook_state) -> string
    + returns a JSON document containing the ordered list of cells with outputs
    # export
    -> std.json.encode_object
