# Requirement: "a terminal session manager that picks a workspace and opens a named session"

Lists candidate workspaces ranked by a recency database, then creates or reattaches to a named session for the chosen workspace.

std
  std.fs
    std.fs.exists
      fn (path: string) -> bool
      + returns true when a path exists
      # filesystem
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns direct children of a directory
      - returns error when the path is not a directory
      # filesystem
  std.proc
    std.proc.run
      fn (cmd: string, args: list[string]) -> result[i32, string]
      + runs a subprocess and returns its exit status
      - returns error when the binary cannot be launched
      # process
    std.proc.run_capture
      fn (cmd: string, args: list[string]) -> result[string, string]
      + runs a subprocess and returns its stdout
      - returns error on non-zero exit
      # process

sessions
  sessions.new
    fn (session_tool: string) -> manager_state
    + creates a manager bound to a session-multiplexer command name
    # construction
  sessions.load_workspaces
    fn (state: manager_state, roots: list[string]) -> result[manager_state, string]
    + scans each root directory one level deep and records each child as a candidate workspace
    - returns error when a root cannot be listed
    # discovery
    -> std.fs.list_dir
  sessions.load_recency
    fn (state: manager_state, frecency_cmd: string) -> result[manager_state, string]
    + asks an external frecency database for scored paths and attaches scores to workspaces
    - returns error when the command fails
    # recency
    -> std.proc.run_capture
  sessions.list
    fn (state: manager_state) -> list[workspace_entry]
    + returns workspaces sorted by recency score descending, then name ascending
    # listing
  sessions.session_name
    fn (path: string) -> string
    + derives a session name from a workspace path by taking the final segment
    # naming
  sessions.open
    fn (state: manager_state, workspace: string) -> result[void, string]
    + attaches to an existing session for the workspace, creating one if none exists
    - returns error when the workspace is not found
    # open
    -> std.fs.exists
    -> std.proc.run
  sessions.kill
    fn (state: manager_state, name: string) -> result[void, string]
    + terminates the named session
    # kill
    -> std.proc.run
