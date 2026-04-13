# Requirement: "a dynamic decision tree where callers define rules and evaluate them against facts"

Rules are expressions over a fact map; leaves carry an outcome label. Trees are built programmatically.

std: (all units exist)

decision
  decision.new_tree
    @ () -> tree_state
    + creates a tree with a single empty root node
    # construction
  decision.add_branch
    @ (state: tree_state, parent_path: list[i32], condition: string) -> result[tuple[list[i32], tree_state], string]
    + appends a new child under the parent with a condition expression
    + returns the new child's path
    - returns error when parent_path does not resolve to a node
    - returns error when condition fails to parse
    # tree_building
  decision.set_outcome
    @ (state: tree_state, path: list[i32], label: string) -> result[tree_state, string]
    + marks a node as a leaf with the given outcome label
    - returns error when the path does not resolve
    # tree_building
  decision.evaluate
    @ (state: tree_state, facts: map[string, string]) -> result[string, string]
    + walks the tree choosing the first branch whose condition matches the facts
    + returns the outcome label of the reached leaf
    - returns error when no branch matches at some interior node
    - returns error when no leaf has an outcome set
    # evaluation
  decision.trace
    @ (state: tree_state, facts: map[string, string]) -> result[list[string], string]
    + returns the sequence of condition strings taken to reach a leaf
    - returns error when evaluation fails
    # introspection
  decision.parse_condition
    @ (source: string) -> result[condition_expr, string]
    + parses a simple boolean expression over fact variables
    - returns error on unknown operator
    # parsing
