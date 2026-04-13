# Requirement: "a terminal UI playground for experimenting with text-processing commands"

The user types a command template and sample input; the library runs the command through a pluggable executor and renders output as a frame. All execution goes through an injected executor so tests are deterministic.

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

playground
  playground.new
    @ (executor: command_executor) -> playground_state
    + creates a playground bound to the given executor
    # construction
  playground.set_command
    @ (state: playground_state, command: string) -> playground_state
    + updates the command template
    # editing
  playground.set_input
    @ (state: playground_state, input: string) -> playground_state
    + updates the sample input
    # editing
  playground.run
    @ (state: playground_state) -> playground_state
    + runs command with input via the executor and stores stdout, stderr, exit
    - stores an error record when the executor rejects the command
    # execution
  playground.history
    @ (state: playground_state) -> list[run_record]
    + returns all prior runs in chronological order
    # history
  playground.undo
    @ (state: playground_state) -> playground_state
    + restores command and input to the previous history entry
    # history
  playground.render
    @ (state: playground_state, width: i32, height: i32) -> string
    + renders a three-pane view: command, input, output
    # rendering
    -> std.tui.new_screen
    -> std.tui.draw_text
    -> std.tui.render
