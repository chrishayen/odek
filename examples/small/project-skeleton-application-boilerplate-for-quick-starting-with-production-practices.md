# Requirement: "an application boilerplate for quick-starting projects following production practices"

A small library that produces a project skeleton description from a name and a set of feature flags.

std: (all units exist)

project_skeleton
  project_skeleton.new
    fn (name: string) -> project_skeleton_state
    + creates a skeleton with the given project name and no features enabled
    - treats an empty name as an invalid skeleton
    # construction
  project_skeleton.enable_feature
    fn (state: project_skeleton_state, feature: string) -> project_skeleton_state
    + records a feature flag on the skeleton
    # configuration
  project_skeleton.file_plan
    fn (state: project_skeleton_state) -> list[string]
    + returns the ordered list of files the skeleton would create
    + includes feature-specific files only when those features are enabled
    # planning
