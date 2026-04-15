# Requirement: "a terminal file explorer library with a navigation model decoupled from rendering"

Project state tracks current directory, selection, and a filter. Rendering is the caller's responsibility.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[dir_entry], string]
      + returns entries of a directory with name and is_dir flag
      - returns error when the directory does not exist
      # filesystem
    std.fs.parent
      fn (path: string) -> string
      + returns the parent directory path
      + returns the same path at the root
      # filesystem
    std.fs.join
      fn (base: string, name: string) -> string
      + joins a base directory and a child name
      # filesystem

file_explorer
  file_explorer.open
    fn (path: string) -> result[explorer_state, string]
    + initializes the explorer at the given directory with the first entry selected
    - returns error when the directory cannot be read
    # construction
    -> std.fs.list_dir
  file_explorer.entries
    fn (state: explorer_state) -> list[dir_entry]
    + returns the entries currently visible, after filter application
    # query
  file_explorer.selected
    fn (state: explorer_state) -> optional[dir_entry]
    + returns the currently selected entry or none if the list is empty
    # query
  file_explorer.move_cursor
    fn (state: explorer_state, delta: i32) -> explorer_state
    + moves the selection by delta, clamped to the visible range
    # navigation
  file_explorer.enter
    fn (state: explorer_state) -> result[explorer_state, string]
    + descends into the selected directory
    - returns error when the selection is not a directory
    # navigation
    -> std.fs.list_dir
    -> std.fs.join
  file_explorer.ascend
    fn (state: explorer_state) -> result[explorer_state, string]
    + moves to the parent directory
    - returns error when already at the filesystem root
    # navigation
    -> std.fs.parent
    -> std.fs.list_dir
  file_explorer.set_filter
    fn (state: explorer_state, pattern: string) -> explorer_state
    + applies a substring filter to the visible entries
    # filtering
