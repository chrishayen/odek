# Requirement: "a curated link list library"

Holds a list of titled links and returns them. The only meaningful operation is appending and listing.

std: (all units exist)

link_list
  link_list.new
    @ () -> link_list_state
    + returns an empty link list
    # construction
  link_list.add
    @ (state: link_list_state, title: string, url: string) -> link_list_state
    + appends a titled link to the list
    # mutation
  link_list.all
    @ (state: link_list_state) -> list[tuple[string, string]]
    + returns all (title, url) pairs in insertion order
    # read
