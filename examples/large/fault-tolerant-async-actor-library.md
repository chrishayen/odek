# Requirement: "a fault-tolerant async actor library"

Actors own state, process one message at a time from a mailbox, and can be supervised so that a panicking child is restarted by its parent.

std
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.async
    std.async.spawn
      @ (task: fn() -> void) -> task_handle
      + schedules task on the async runtime
      # scheduling
    std.async.sleep_ms
      @ (ms: i64) -> void
      + yields the current task for at least ms milliseconds
      # scheduling
  std.sync
    std.sync.channel
      @ (capacity: i32) -> tuple[channel_sender, channel_receiver]
      + creates a bounded async channel
      # synchronization
    std.sync.channel_send
      @ (sender: channel_sender, msg: bytes) -> result[void, string]
      + sends a message, awaiting capacity if full
      - returns error when the receiver is dropped
      # synchronization
    std.sync.channel_recv
      @ (receiver: channel_receiver) -> result[bytes, string]
      + receives the next message
      - returns error when the channel is closed
      # synchronization

actor
  actor.spawn
    @ (initial_state: bytes, handler: fn(bytes, bytes) -> result[bytes, string]) -> actor_ref
    + starts an actor with the given handler and returns a reference
    ? handler receives (state, message) and returns new state or an error
    # lifecycle
    -> std.async.spawn
    -> std.sync.channel
  actor.send
    @ (ref: actor_ref, message: bytes) -> result[void, string]
    + enqueues a message into the actor's mailbox
    - returns error when the actor has stopped
    # messaging
    -> std.sync.channel_send
  actor.ask
    @ (ref: actor_ref, message: bytes, timeout_ms: i64) -> result[bytes, string]
    + sends a message and waits for the actor's reply
    - returns error when the timeout elapses
    # messaging
    -> std.time.now_millis
    -> std.async.sleep_ms
  actor.stop
    @ (ref: actor_ref) -> void
    + signals the actor to drain and terminate
    # lifecycle
  actor.supervise
    @ (parent: actor_ref, child: actor_ref, strategy: restart_strategy) -> void
    + links parent to child so parent is notified on child failure
    # supervision
  actor.on_failure
    @ (strategy: restart_strategy, child: actor_ref, error: string) -> supervisor_decision
    + decides whether to restart, stop, or escalate the failing child
    + one_for_one restarts only the failing child
    + one_for_all restarts every sibling
    # supervision
  actor.restart
    @ (ref: actor_ref) -> result[void, string]
    + drains the mailbox and reinitializes the actor state
    - returns error when the actor has been stopped permanently
    # supervision
    -> std.async.spawn
  actor.run_loop
    @ (ref: actor_ref) -> void
    + reads the next message and invokes the handler until the actor is stopped
    # runtime
    -> std.sync.channel_recv
