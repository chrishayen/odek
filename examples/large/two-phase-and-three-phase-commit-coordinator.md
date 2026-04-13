# Requirement: "a distributed transactions coordinator supporting two-phase and three-phase commit"

The coordinator tracks transaction state across multiple participants. Participants are addressed by opaque id; the transport is the caller's responsibility.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.new_uuid
      @ () -> string
      + returns a fresh random identifier
      # identifiers

txn
  txn.new_coordinator
    @ () -> coordinator_state
    + creates a coordinator with no open transactions
    # construction
  txn.begin
    @ (state: coordinator_state, participants: list[string]) -> tuple[string, coordinator_state]
    + allocates a transaction id and records the participant set in state INIT
    # transaction_lifecycle
    -> std.id.new_uuid
    -> std.time.now_millis
  txn.record_vote
    @ (state: coordinator_state, txn_id: string, participant: string, yes: bool) -> result[coordinator_state, string]
    + records a participant's vote during the voting phase
    - returns error when the transaction is not in VOTING state
    - returns error when the participant was not enrolled
    # two_phase_commit
  txn.decide_2pc
    @ (state: coordinator_state, txn_id: string) -> result[tuple[string, coordinator_state], string]
    + returns "commit" when all participants voted yes and transitions to COMMITTED
    + returns "abort" when any participant voted no and transitions to ABORTED
    - returns error when not all votes have been recorded
    # two_phase_commit
  txn.decide_3pc_precommit
    @ (state: coordinator_state, txn_id: string) -> result[coordinator_state, string]
    + transitions from VOTING to PRECOMMIT when all votes are yes
    - returns error when any vote is no
    - returns error when votes are incomplete
    # three_phase_commit
  txn.decide_3pc_commit
    @ (state: coordinator_state, txn_id: string) -> result[coordinator_state, string]
    + transitions from PRECOMMIT to COMMITTED
    - returns error when the transaction is not in PRECOMMIT state
    # three_phase_commit
  txn.abort
    @ (state: coordinator_state, txn_id: string) -> result[coordinator_state, string]
    + moves a transaction to ABORTED from any non-terminal state
    - returns error when the transaction id is unknown
    # transaction_lifecycle
  txn.status
    @ (state: coordinator_state, txn_id: string) -> result[string, string]
    + returns the current state name for the transaction
    - returns error when the transaction id is unknown
    # introspection
  txn.timeout_sweep
    @ (state: coordinator_state, stale_ms: i64) -> coordinator_state
    + aborts any transaction whose last update is older than stale_ms
    + leaves active transactions unchanged
    # recovery
    -> std.time.now_millis
