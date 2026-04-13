# Requirement: "a terminal UI for running and monitoring multiple processes"

Process spawning is injected. The library owns the process-list state machine and renders a tab bar with an output pane.

std
  std.tui
    std.tui.new_screen
      @ (width: i32, height: i32) -> screen
      + creates an off-screen character buffer
      # tui_primitive
    std.tui.draw_text
      @ (s: screen, row: i32, col: i32, text: string) -> screen
      + writes text at the given position
      # tui_primitive
    std.tui.render
      @ (s: screen) -> string
      + returns the screen as a newline-delimited string
      # tui_primitive

mprocs
  mprocs.new
    @ (spawner: process_spawner) -> mprocs_state
    + creates an empty process list bound to a spawner
    # construction
  mprocs.add
    @ (state: mprocs_state, name: string, command: string) -> mprocs_state
    + registers a process entry without starting it
    # registration
  mprocs.start
    @ (state: mprocs_state, index: i32) -> result[mprocs_state, string]
    + launches the process at index via the spawner
    - returns error when the process is already running
    # lifecycle
  mprocs.stop
    @ (state: mprocs_state, index: i32) -> result[mprocs_state, string]
    + signals the process at index to terminate
    - returns error when the process is not running
    # lifecycle
  mprocs.restart
    @ (state: mprocs_state, index: i32) -> result[mprocs_state, string]
    + stops then starts the process at index
    # lifecycle
  mprocs.append_output
    @ (state: mprocs_state, index: i32, line: string) -> mprocs_state
    + appends a line to the process's output buffer, capping at max lines
    # output_buffer
  mprocs.select
    @ (state: mprocs_state, index: i32) -> mprocs_state
    + switches the active tab to the given process
    # ui_state
  mprocs.render
    @ (state: mprocs_state, width: i32, height: i32) -> string
    + renders a tab bar above the selected process output
    # rendering
    -> std.tui.new_screen
    -> std.tui.draw_text
    -> std.tui.render
