# Requirement: "a disk usage analyzer with a console interface"

Walks a directory tree, computes cumulative sizes, and presents a navigable tree view in the terminal.

std
  std.fs
    std.fs.stat
      @ (path: string) -> result[file_info, string]
      + returns size, modified time, and kind (file/dir/symlink)
      - returns error when path does not exist
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns the names of directory entries
      - returns error on permission denied
      # filesystem
  std.term
    std.term.raw_mode_enter
      @ () -> result[term_state, string]
      + switches terminal to raw mode and returns a restore handle
      # terminal
    std.term.raw_mode_leave
      @ (state: term_state) -> void
      + restores prior terminal mode
      # terminal
    std.term.read_key
      @ () -> result[key_event, string]
      + blocks for a key press
      # terminal
    std.term.draw
      @ (frame: list[string]) -> result[void, string]
      + replaces the current screen with the given lines
      # terminal
    std.term.size
      @ () -> tuple[i32, i32]
      + returns terminal (cols, rows)
      # terminal

disk_usage
  disk_usage.scan
    @ (root: string) -> result[disk_tree, string]
    + walks root recursively and builds a tree with per-node cumulative size
    ? symlinks are not followed; directories that error are recorded and skipped
    # scanning
    -> std.fs.stat
    -> std.fs.list_dir
  disk_usage.sort_children_by_size
    @ (tree: disk_tree) -> disk_tree
    + sorts children at every level largest-first
    # sorting
  disk_usage.render_frame
    @ (tree: disk_tree, cursor: list[string], cols: i32, rows: i32) -> list[string]
    + returns rendered lines for the current cursor path, including size bars and percentages
    # rendering
  disk_usage.handle_key
    @ (tree: disk_tree, cursor: list[string], key: key_event) -> list[string]
    + returns the updated cursor after applying a navigation key
    + supports up/down to move between siblings and enter/backspace to descend and ascend
    # navigation
  disk_usage.run_tui
    @ (root: string) -> result[void, string]
    + scans root then runs an interactive terminal ui until the user quits
    # orchestration
    -> std.term.raw_mode_enter
    -> std.term.raw_mode_leave
    -> std.term.read_key
    -> std.term.draw
    -> std.term.size
