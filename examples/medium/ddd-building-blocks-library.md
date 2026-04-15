# Requirement: "a building-block library for domain-driven design"

Primitives for aggregates, domain events, and a unit-of-work envelope that application services compose into use cases.

std: (all units exist)

ddd
  ddd.new_aggregate
    fn (id: string, version: i64) -> aggregate_state
    + creates an aggregate root with the given id and initial version
    # construction
  ddd.record_event
    fn (agg: aggregate_state, event_type: string, payload: bytes) -> aggregate_state
    + appends a domain event to the aggregate's uncommitted queue and bumps the pending version
    # events
  ddd.uncommitted_events
    fn (agg: aggregate_state) -> list[domain_event]
    + returns the events recorded since the last mark_committed
    # events
  ddd.mark_committed
    fn (agg: aggregate_state) -> aggregate_state
    + clears the uncommitted queue and promotes the pending version
    # events
  ddd.new_unit_of_work
    fn () -> uow_state
    + creates an empty unit-of-work collecting aggregates to persist together
    # construction
  ddd.register
    fn (uow: uow_state, agg: aggregate_state) -> uow_state
    + adds an aggregate to the unit of work
    - replaces a prior registration with the same id
    # transactions
  ddd.commit
    fn (uow: uow_state, persist: fn(aggregate_state) -> result[void, string], publish: fn(domain_event) -> void) -> result[void, string]
    + persists every registered aggregate, then publishes each uncommitted event
    - aborts on the first persist failure and publishes nothing
    # transactions
    -> ddd.uncommitted_events
    -> ddd.mark_committed
  ddd.command_handler
    fn (uow: uow_state, decide: fn(uow_state) -> result[uow_state, string]) -> result[uow_state, string]
    + runs a pure decision function that may register or mutate aggregates
    - propagates any error from decide without mutating the uow
    # application
