# Requirement: "a library for updating terminal output in real time"

Stateful writer that overwrites the last frame by emitting cursor-up and line-clear escape sequences before each new frame.

std
  std.io
    std.io.write_stdout
      fn (data: string) -> result[void, string]
      + writes the string to standard output
      - returns error when stdout is closed
      # io

live_writer
  live_writer.new
    fn () -> live_writer_state
    + creates a writer with no previously-rendered frame
    # construction
  live_writer.render
    fn (state: live_writer_state, frame: string) -> result[live_writer_state, string]
    + emits escape sequences to erase the prior frame then writes the new frame
    + on the first call nothing is erased
    ? counts lines in the previous frame to decide how many rows to clear
    # rendering
    -> std.io.write_stdout
  live_writer.finish
    fn (state: live_writer_state) -> result[void, string]
    + writes a trailing newline so subsequent output starts on a fresh line
    # teardown
    -> std.io.write_stdout
