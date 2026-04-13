# Requirement: "a library that drives a debugger via its machine interface protocol"

Speaks the GDB/MI protocol: sends commands, parses asynchronous output records, and exposes breakpoint and stepping operations.

std
  std.process
    std.process.spawn
      @ (command: string, args: list[string]) -> result[process_handle, string]
      + launches the subprocess and returns a handle
      - returns error when the executable cannot be found
      # process
    std.process.write_stdin
      @ (handle: process_handle, data: string) -> result[void, string]
      + writes a line to the process stdin
      # process
    std.process.read_stdout_line
      @ (handle: process_handle) -> result[string, string]
      + reads one line from the process stdout
      # process

debugger
  debugger.attach
    @ (target: string) -> result[debugger_state, string]
    + spawns a debugger attached to the target binary
    - returns error when the target is not found
    # connection
    -> std.process.spawn
  debugger.send_command
    @ (state: debugger_state, command: string) -> result[debugger_state, string]
    + sends a machine-interface command with a token prefix
    # protocol
    -> std.process.write_stdin
  debugger.parse_record
    @ (line: string) -> result[debugger_record, string]
    + parses one MI result, async, or stream record
    - returns error on malformed records
    # protocol
  debugger.set_breakpoint
    @ (state: debugger_state, location: string) -> result[i32, string]
    + returns the new breakpoint number
    - returns error when the location cannot be resolved
    # breakpoints
  debugger.step
    @ (state: debugger_state) -> result[debugger_record, string]
    + advances execution one source line
    - returns error when not stopped
    # stepping
    -> std.process.read_stdout_line
  debugger.read_stack
    @ (state: debugger_state) -> result[list[stack_frame], string]
    + returns the current call stack
    # inspection
