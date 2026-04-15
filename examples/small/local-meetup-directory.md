# Requirement: "a local meetup directory"

The source entry was a meetup listing rather than a library idea. Interpreted as a tiny catalog for meetup groups by city.

std: (all units exist)

meetups
  meetups.new
    fn () -> directory_state
    + returns an empty meetup directory
    # construction
  meetups.add_group
    fn (dir: directory_state, name: string, city: string, url: string) -> directory_state
    + registers a meetup group under its city
    # registration
  meetups.list_by_city
    fn (dir: directory_state, city: string) -> list[group]
    + returns every group registered in city, case-insensitive match
    + returns an empty list when the city has no groups
    # query
  meetups.find
    fn (dir: directory_state, name: string) -> optional[group]
    + returns the group with the given name, if any
    # query
