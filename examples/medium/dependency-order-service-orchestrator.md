# Requirement: "a simple service orchestrator that starts services in dependency order"

Given a set of named services with dependencies, produces a start order, starts each through a pluggable runner, and tracks health.

std: (all units exist)

orchestrator
  orchestrator.new
    @ () -> orchestrator_state
    + constructs an empty orchestrator with no services
    # construction
  orchestrator.register
    @ (state: orchestrator_state, name: string, depends_on: list[string], runner: service_runner) -> result[orchestrator_state, string]
    + returns a new state with the service added
    - returns error when the name is already registered
    # registration
  orchestrator.resolve_order
    @ (state: orchestrator_state) -> result[list[string], string]
    + returns names in topological order so each service sees its deps started first
    - returns error when a dependency is missing
    - returns error when a cycle is detected
    # ordering
  orchestrator.start_all
    @ (state: orchestrator_state) -> result[orchestrator_state, string]
    + starts every service via its runner in resolved order
    - stops on first failure and returns error
    # startup
  orchestrator.stop_all
    @ (state: orchestrator_state) -> result[orchestrator_state, string]
    + stops services in reverse start order
    - continues past individual failures and returns the first error at the end
    # shutdown
  orchestrator.health
    @ (state: orchestrator_state, name: string) -> result[string, string]
    + returns "starting", "running", "stopped", or "failed"
    - returns error when no such service is registered
    # health
