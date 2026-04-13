# Requirement: "a terminal screensaver that visualizes a commit history"

Takes a list of commits and produces animation frames. Repository loading is a separate std utility so the visualizer can be fed synthetic data in tests.

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
  std.vcs
    std.vcs.load_commits
      @ (repo_path: string, limit: i32) -> result[list[commit], string]
      + returns up to limit commits from HEAD in reverse chronological order
      - returns error when repo_path is not a repository
      # version_control

commit_screensaver
  commit_screensaver.new
    @ (commits: list[commit]) -> screensaver_state
    + creates a screensaver seeded with the commit list
    # construction
  commit_screensaver.load_from_repo
    @ (repo_path: string, limit: i32) -> result[screensaver_state, string]
    + loads commits via the std helper and constructs a state
    # loading
    -> std.vcs.load_commits
  commit_screensaver.advance
    @ (state: screensaver_state) -> screensaver_state
    + advances the animation by one tick
    # animation
  commit_screensaver.layout_graph
    @ (commits: list[commit], width: i32) -> list[graph_row]
    + computes column placement for each commit so parent lines do not cross
    # graph_layout
  commit_screensaver.render
    @ (state: screensaver_state, width: i32, height: i32) -> string
    + renders the current animation frame
    # rendering
    -> std.tui.new_screen
    -> std.tui.draw_text
    -> std.tui.render
