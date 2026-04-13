# Requirement: "a fuzzy-search terminal selector library for issue trackers"

A generic fuzzy selector over issues fetched through a pluggable tracker adapter. The caller renders and handles input.

std
  std.term
    std.term.clear_screen
      @ () -> string
      + returns the ANSI sequence that clears the terminal
      # terminal
    std.term.move_cursor
      @ (row: i32, col: i32) -> string
      + returns the ANSI sequence that moves the cursor
      # terminal

issue_picker
  issue_picker.new
    @ (adapter: tracker_adapter) -> picker_state
    + creates a picker backed by the given tracker adapter with an empty query
    # construction
  issue_picker.load_issues
    @ (state: picker_state, project: string) -> result[picker_state, string]
    + fetches issues for the project via the adapter and stores them in state
    - returns error when the adapter fails
    # data
  issue_picker.update_query
    @ (state: picker_state, query: string) -> picker_state
    + sets the current query and recomputes ranked matches
    + an empty query ranks issues by recency
    # query
    -> issue_picker.score_match
  issue_picker.score_match
    @ (candidate: string, query: string) -> i32
    + returns a fuzzy score where higher means closer
    + returns 0 when no characters of query appear in order in candidate
    # ranking
  issue_picker.current_matches
    @ (state: picker_state) -> list[issue_record]
    + returns the currently visible issues in rank order
    # query
  issue_picker.move_selection
    @ (state: picker_state, delta: i32) -> picker_state
    + moves the highlighted row by delta, clamped to the visible list
    # ui
  issue_picker.selected_issue
    @ (state: picker_state) -> optional[issue_record]
    + returns the currently highlighted issue
    - returns none when there are no matches
    # ui
  issue_picker.render
    @ (state: picker_state) -> string
    + returns a full-screen string showing the query and top matches
    # ui
    -> std.term.clear_screen
    -> std.term.move_cursor
