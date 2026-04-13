# Requirement: "a library that picks this week's discussion project from a rotating roster"

Given a roster and a week number, return the project due this week.

std: (all units exist)

roster_pick
  roster_pick.select_for_week
    @ (roster: list[string], week_index: i64) -> optional[string]
    + returns the roster entry at week_index modulo roster length
    - returns none when the roster is empty
    # selection
