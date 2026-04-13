# Requirement: "a fluent state machine library"

Define states, transitions, and guards with a builder, then execute transitions against an instance.

std: (all units exist)

fsm
  fsm.new
    @ (initial: string) -> machine_definition
    + creates a machine definition with an initial state
    # construction
  fsm.add_state
    @ (def: machine_definition, name: string) -> machine_definition
    + registers a state by name
    ? re-registering an existing state is a no-op
    # definition
  fsm.on
    @ (def: machine_definition, from: string, event: string, to: string) -> machine_definition
    + registers a transition from "from" to "to" triggered by "event"
    # definition
  fsm.guard
    @ (def: machine_definition, from: string, event: string, g: guard_fn) -> machine_definition
    + attaches a guard function that must return true for the transition to fire
    # definition
  fsm.validate
    @ (def: machine_definition) -> result[void, string]
    + checks that every transition references a registered state
    - returns error naming the first missing state
    # validation
  fsm.instance
    @ (def: machine_definition) -> machine_instance
    + creates a runtime instance positioned at the initial state
    # runtime
  fsm.current
    @ (inst: machine_instance) -> string
    + returns the current state name
    # runtime
  fsm.fire
    @ (inst: machine_instance, event: string, context: map[string, string]) -> result[machine_instance, string]
    + applies the matching transition and returns a new instance at the target state
    - returns error when no transition exists for (current, event)
    - returns error when the guard rejects the transition
    # runtime
  fsm.can_fire
    @ (inst: machine_instance, event: string) -> bool
    + returns true when a transition exists for (current, event) ignoring guards
    # introspection
