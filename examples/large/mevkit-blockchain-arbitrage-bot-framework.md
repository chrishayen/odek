# Requirement: "a modular framework for writing blockchain arbitrage bots"

An event-driven engine that wires pluggable collectors (event sources), strategies (decision makers), and executors (action sinks). Collectors and executors are opaque handles supplied by the caller.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.async
    std.async.channel
      fn (capacity: i32) -> async_channel
      + creates a bounded asynchronous channel
      # async
    std.async.send
      fn (ch: async_channel, item: bytes) -> result[void, string]
      + blocks until the item is accepted
      - returns error when the channel is closed
      # async
    std.async.recv
      fn (ch: async_channel) -> result[optional[bytes], string]
      + returns the next item or none when the channel is closed and empty
      # async
    std.async.spawn
      fn (task: async_task) -> async_handle
      + schedules a task to run concurrently
      # async
  std.log
    std.log.info
      fn (message: string) -> void
      + emits an informational log entry
      # logging

mevkit
  mevkit.new_engine
    fn () -> engine_state
    + creates an empty engine with no collectors, strategies, or executors
    # construction
  mevkit.add_collector
    fn (eng: engine_state, name: string, collector: collector_handle) -> engine_state
    + registers an event source under a unique name
    # wiring
  mevkit.add_strategy
    fn (eng: engine_state, name: string, strategy: strategy_handle) -> engine_state
    + registers a decision maker under a unique name
    # wiring
  mevkit.add_executor
    fn (eng: engine_state, name: string, executor: executor_handle) -> engine_state
    + registers an action sink under a unique name
    # wiring
  mevkit.collect_once
    fn (eng: engine_state, collector_name: string) -> result[list[bytes], string]
    + drains currently available events from the named collector
    - returns error when collector_name is unknown
    # collection
    -> std.time.now_millis
  mevkit.dispatch
    fn (eng: engine_state, event: bytes) -> result[list[bytes], string]
    + feeds an event through every strategy and returns the resulting action list
    # strategy
    -> std.log.info
  mevkit.execute
    fn (eng: engine_state, executor_name: string, action: bytes) -> result[void, string]
    + submits an action to the named executor
    - returns error when executor_name is unknown
    # execution
  mevkit.run
    fn (eng: engine_state) -> result[void, string]
    + runs collect/dispatch/execute in a loop until all collectors close
    - returns error on the first unrecoverable executor failure
    # orchestration
    -> std.async.channel
    -> std.async.send
    -> std.async.recv
    -> std.async.spawn
  mevkit.stop
    fn (eng: engine_state) -> result[void, string]
    + signals every registered component to shut down
    # lifecycle
  mevkit.stats
    fn (eng: engine_state) -> engine_stats
    + returns counters for events observed, actions produced, and actions executed
    # introspection
