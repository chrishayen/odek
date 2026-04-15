# Requirement: "a recommended project layout generator"

The input is a project name; the output is a list of recommended directories and their purposes. This is content, not logic.

std: (all units exist)

layout
  layout.standard_directories
    fn () -> list[directory_entry]
    + returns the canonical list of top-level directories and their intended purpose
    ? entries are hardcoded; not configurable
    # layout
  layout.render
    fn (project_name: string) -> string
    + returns a printable tree diagram showing the layout rooted at project_name
    # rendering
