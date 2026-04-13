# Requirement: "a terminal UI for browsing, comparing, and editing dotenv files"

Parsing, diffing, and editing are pure functions over in-memory env maps. Rendering produces a frame string.

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
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the entire file contents
      - returns error when the path does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + writes content, truncating any previous file
      # filesystem

env_tui
  env_tui.parse
    @ (text: string) -> list[env_entry]
    + parses KEY=VALUE lines, preserving order and comments
    + strips surrounding single or double quotes from values
    ? lines starting with # are retained as comment entries
    # parsing
  env_tui.serialize
    @ (entries: list[env_entry]) -> string
    + renders entries back to dotenv text in original order
    # serialization
  env_tui.diff
    @ (left: list[env_entry], right: list[env_entry]) -> list[env_diff_row]
    + returns one row per key with left-only, right-only, same, or changed
    # diff
  env_tui.set
    @ (entries: list[env_entry], key: string, value: string) -> list[env_entry]
    + updates an existing key in place or appends a new one
    # editing
  env_tui.remove
    @ (entries: list[env_entry], key: string) -> list[env_entry]
    + removes the entry with the given key if present
    # editing
  env_tui.load_pair
    @ (left_path: string, right_path: string) -> result[tuple[list[env_entry], list[env_entry]], string]
    + loads and parses both files for side-by-side compare
    - returns error when either file cannot be read
    # loading
    -> std.fs.read_all
  env_tui.save
    @ (path: string, entries: list[env_entry]) -> result[void, string]
    + writes serialized entries to disk
    # saving
    -> std.fs.write_all
  env_tui.render
    @ (rows: list[env_diff_row], selected: i32, width: i32, height: i32) -> string
    + renders a two-column diff view with the selected row highlighted
    # rendering
    -> std.tui.new_screen
    -> std.tui.draw_text
    -> std.tui.render
