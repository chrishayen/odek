# Requirement: "a community-wiki project list library"

Holds projects with title and description and returns them. Insertion and listing only.

std: (all units exist)

projects
  projects.new
    fn () -> projects_state
    + returns an empty list
    # construction
  projects.add
    fn (state: projects_state, title: string, description: string, url: string) -> projects_state
    + appends a project entry
    # mutation
  projects.all
    fn (state: projects_state) -> list[tuple[string, string, string]]
    + returns all (title, description, url) triples in insertion order
    # read
