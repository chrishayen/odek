# Requirement: "a library for live presentations in a terminal"

Plays back a script of shell commands one keypress at a time, as if the presenter were typing them live.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when the file cannot be opened
      # filesystem
  std.io
    std.io.read_key
      @ () -> result[string, string]
      + returns a single keystroke from the input stream
      - returns error when the stream is closed
      # input

live_present
  live_present.load_script
    @ (path: string) -> result[list[string], string]
    + returns one command per non-empty, non-comment line of the file
    - returns error when the file cannot be read
    # loading
    -> std.fs.read_all
  live_present.type_command
    @ (command: string) -> result[void, string]
    + emits the command one character at a time, advancing on each keypress from the input stream
    ? a final keypress is required before the command is considered typed
    # playback
    -> std.io.read_key
  live_present.run_script
    @ (commands: list[string]) -> result[void, string]
    + walks the commands in order, typing each one and waiting for confirmation before yielding it to the caller for execution
    - returns error when the input stream closes mid-script
    # playback
