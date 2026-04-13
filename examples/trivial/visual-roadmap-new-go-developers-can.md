# Requirement: "a learning roadmap step tracker"

Represents an ordered sequence of learning steps and tracks which have been completed.

std: (all units exist)

roadmap
  roadmap.new
    @ (steps: list[string]) -> roadmap_state
    + creates a roadmap with the given steps, none completed
    # construction
  roadmap.mark_complete
    @ (state: roadmap_state, step: string) -> roadmap_state
    + marks the named step as completed
    - is a no-op when the step is not in the roadmap
    # writes
  roadmap.progress
    @ (state: roadmap_state) -> tuple[i32, i32]
    + returns (completed_count, total_count)
    # reads
