# Requirement: "a structured learning path library"

Represents a tree of learning nodes with prerequisites and tracks which nodes a learner has completed, exposing what is unlocked next.

std: (all units exist)

learnpath
  learnpath.new_tree
    @ () -> tree_state
    + creates an empty learning path tree
    # construction
  learnpath.add_node
    @ (tree: tree_state, node_id: string, title: string, prerequisites: list[string]) -> result[tree_state, string]
    + adds a node with its prerequisite node ids
    - returns error when a prerequisite refers to an unknown node
    - returns error when node_id already exists
    # schema
  learnpath.mark_completed
    @ (tree: tree_state, learner_id: string, node_id: string) -> result[tree_state, string]
    + records that learner completed node_id
    - returns error when node_id does not exist
    # progress
  learnpath.available_nodes
    @ (tree: tree_state, learner_id: string) -> list[string]
    + returns node ids whose prerequisites are all completed and which are not yet done
    # query
  learnpath.progress_ratio
    @ (tree: tree_state, learner_id: string) -> f64
    + returns fraction of total nodes the learner has completed (0.0 to 1.0)
    # query
