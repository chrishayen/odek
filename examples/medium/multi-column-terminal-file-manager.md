# Requirement: "a library for a multi-column terminal file manager"

Stateful three-pane browser: parent, current, preview. Library exposes state plus navigation and selection commands. Rendering to a terminal is the caller's job.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[dir_entry], string]
      + returns entries with name, is_dir, and size
      - returns error when the path is not a directory
      # filesystem
    std.fs.read_head
      fn (path: string, max_bytes: i64) -> result[bytes, string]
      + reads up to max_bytes from the start of a file
      - returns error on permission denied
      # filesystem

file_manager
  file_manager.new
    fn (root: string) -> result[fm_state, string]
    + returns initial state with cursor at the first entry of root
    - returns error when root is not a directory
    # construction
    -> std.fs.list_dir
  file_manager.move_cursor
    fn (state: fm_state, delta: i32) -> fm_state
    + moves the cursor by delta, clamped to the entry list
    # navigation
  file_manager.enter
    fn (state: fm_state) -> result[fm_state, string]
    + enters the directory under the cursor
    - returns error when the cursor is on a file
    # navigation
    -> std.fs.list_dir
  file_manager.go_up
    fn (state: fm_state) -> result[fm_state, string]
    + moves to the parent directory and places the cursor on the previous child
    - returns error at filesystem root
    # navigation
    -> std.fs.list_dir
  file_manager.toggle_select
    fn (state: fm_state) -> fm_state
    + toggles selection on the cursor entry
    # selection
  file_manager.preview
    fn (state: fm_state) -> optional[bytes]
    + returns a preview head of the cursor entry when it is a file, none otherwise
    # preview
    -> std.fs.read_head
