# Requirement: "a graph-based control-flow state machine with aspect hooks"

Declare named states with relationships (requires, removes, after), transition by asserting state sets, and fire before/after aspect callbacks around each transition.

std: (all units exist)

flow_machine
  flow_machine.new
    @ () -> machine_state
    + creates an empty machine with no states declared
    # construction
  flow_machine.declare_state
    @ (state: machine_state, name: string, requires: list[string], removes: list[string]) -> result[machine_state, string]
    + registers a state along with the states it requires and removes on entry
    - returns error when name is already declared
    # schema
  flow_machine.declare_order
    @ (state: machine_state, before: string, after: string) -> result[machine_state, string]
    + declares that "before" must be asserted strictly prior to "after" within a transition
    - returns error when either state is unknown
    # schema
  flow_machine.add_hook
    @ (state: machine_state, target: string, phase: string, hook: closure[void]) -> result[machine_state, string]
    + registers a hook to run on "before" or "after" entering the target state
    - returns error on unknown phase
    # aspects
  flow_machine.add
    @ (state: machine_state, states: list[string]) -> result[machine_state, string]
    + asserts the listed states, resolving requires/removes and firing hooks
    - returns error when a required state is not satisfied
    - returns error on ordering violation
    # transition
  flow_machine.remove
    @ (state: machine_state, states: list[string]) -> machine_state
    + removes the listed states and any states that depended on them
    # transition
  flow_machine.is_active
    @ (state: machine_state, name: string) -> bool
    + returns true when the named state is currently asserted
    # query
  flow_machine.snapshot
    @ (state: machine_state) -> list[string]
    + returns the set of currently asserted state names
    # query
