# Requirement: "a finite state machine implementation"

Defines states and transitions, then drives a machine forward by dispatching events.

std: (all units exist)

fsm
  fsm.new
    @ (initial_state: string) -> fsm_state
    + returns a machine positioned in the initial state
    # construction
  fsm.add_transition
    @ (machine: fsm_state, from_state: string, event: string, to_state: string) -> fsm_state
    + registers a transition from a state on an event
    # definition
  fsm.dispatch
    @ (machine: fsm_state, event: string) -> result[fsm_state, string]
    + returns a new machine in the target state of the matching transition
    - returns error when no transition is defined for the current state and event
    # dispatch
  fsm.current
    @ (machine: fsm_state) -> string
    + returns the name of the current state
    # inspection
