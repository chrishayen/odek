# Requirement: "a voted resource list library"

Tracks resources with an integer vote count and returns them sorted by votes.

std: (all units exist)

voted_resources
  voted_resources.new
    @ () -> resources_state
    + returns an empty list
    # construction
  voted_resources.add
    @ (state: resources_state, title: string, url: string) -> resources_state
    + appends a resource with zero votes
    # mutation
  voted_resources.upvote
    @ (state: resources_state, title: string) -> resources_state
    + increments the vote count for a matching title
    ? unknown titles are ignored
    # mutation
  voted_resources.ranked
    @ (state: resources_state) -> list[tuple[string, string, i32]]
    + returns (title, url, votes) sorted by votes descending, stable on ties
    # read
