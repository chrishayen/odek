# Requirement: "a BDD-style bifurcation test framework"

Lets the caller describe scenarios as given/when/then steps with a branching tree of preconditions.

std: (all units exist)

biff
  biff.scenario
    @ (name: string) -> scenario_state
    + creates an empty scenario with the given label
    # construction
  biff.given
    @ (state: scenario_state, desc: string, setup: fn() -> context) -> scenario_state
    + adds a precondition node that runs setup to produce a shared context
    # preconditions
  biff.when
    @ (state: scenario_state, desc: string, action: fn(context) -> context) -> scenario_state
    + branches the scenario tree with an action applied to each given context
    # actions
  biff.then
    @ (state: scenario_state, desc: string, check: fn(context) -> bool) -> scenario_state
    + attaches an assertion to every leaf of the current branch
    # assertions
  biff.run
    @ (state: scenario_state) -> list[step_result]
    + walks the tree, runs every given/when/then path, and returns per-branch pass/fail results
    - marks a branch as failed when a then check returns false
    # execution
