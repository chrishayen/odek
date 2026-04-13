# Requirement: "a library for submitting and browsing module-idea requests"

Idea records can be submitted, listed, searched, and upvoted. Purely in-memory state; persistence is the caller's concern.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

module_ideas
  module_ideas.new_board
    @ () -> board_state
    + creates an empty idea board
    # construction
  module_ideas.submit
    @ (state: board_state, title: string, description: string, author: string) -> tuple[string, board_state]
    + adds a new idea and returns its assigned id alongside updated state
    # submission
    -> std.time.now_seconds
  module_ideas.upvote
    @ (state: board_state, id: string, voter: string) -> result[board_state, string]
    + records an upvote from voter on the given idea
    - returns error when the id does not exist
    - returns error when the voter has already upvoted this idea
    # voting
  module_ideas.search
    @ (state: board_state, query: string) -> list[idea]
    + returns ideas whose title or description contains the query, case-insensitive, ranked by upvotes
    # search
  module_ideas.list_recent
    @ (state: board_state, limit: i32) -> list[idea]
    + returns the most recently submitted ideas up to limit
    # listing
