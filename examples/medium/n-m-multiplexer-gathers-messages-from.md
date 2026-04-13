# Requirement: "an n-to-m multiplexer that gathers messages from different sources and broadcasts them to a set of destinations"

Named producers push messages into a shared bus; the bus fans them out to every registered consumer.

std
  std.sync
    std.sync.channel_new
      @ (capacity: i32) -> channel_state
      + creates a bounded channel with the given capacity
      # concurrency
    std.sync.channel_send
      @ (ch: channel_state, payload: bytes) -> result[void, string]
      + enqueues a payload, blocking when the channel is full
      - returns error when the channel is closed
      # concurrency
    std.sync.channel_recv
      @ (ch: channel_state) -> result[bytes, string]
      + dequeues the next payload, blocking until one is available
      - returns error when the channel is closed and drained
      # concurrency

mux_bus
  mux_bus.new
    @ (queue_capacity: i32) -> mux_bus_state
    + creates a bus with no producers, no consumers, and a bounded intake channel
    # construction
    -> std.sync.channel_new
  mux_bus.register_producer
    @ (state: mux_bus_state, name: string) -> tuple[producer_handle, mux_bus_state]
    + assigns an opaque producer handle and records its name
    # registration
  mux_bus.register_consumer
    @ (state: mux_bus_state, name: string) -> tuple[consumer_handle, mux_bus_state]
    + creates a dedicated outbound channel for the consumer
    # registration
    -> std.sync.channel_new
  mux_bus.publish
    @ (state: mux_bus_state, producer: producer_handle, payload: bytes) -> result[void, string]
    + enqueues payload onto the intake channel tagged with producer
    - returns error when producer is not registered
    # ingestion
    -> std.sync.channel_send
  mux_bus.run_dispatch
    @ (state: mux_bus_state) -> result[void, string]
    + drains the intake channel and copies each payload to every registered consumer
    + a blocked consumer does not stall other consumers beyond its own channel
    - returns error when the intake channel is closed and drained
    # dispatch
    -> std.sync.channel_recv
    -> std.sync.channel_send
  mux_bus.consume
    @ (state: mux_bus_state, consumer: consumer_handle) -> result[bytes, string]
    + returns the next payload addressed to the consumer
    - returns error when the consumer's channel is closed
    # consumption
    -> std.sync.channel_recv
  mux_bus.close
    @ (state: mux_bus_state) -> result[void, string]
    + closes intake and all consumer channels so workers can exit
    # lifecycle
