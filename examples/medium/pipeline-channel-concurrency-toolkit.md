# Requirement: "a toolkit for clean, composable, channel-based concurrency"

A small set of pipeline combinators over bounded in-memory channels.

std: (all units exist)

pipeline
  pipeline.new_channel
    fn (capacity: i32) -> channel_state
    + creates a bounded channel with the given capacity
    # construction
  pipeline.send
    fn (ch: channel_state, item: i64) -> bool
    + returns true when the item was accepted
    - returns false when the channel is closed
    # io
  pipeline.receive
    fn (ch: channel_state) -> optional[i64]
    + returns the next item, blocking until one is available
    - returns none when the channel is closed and drained
    # io
  pipeline.close
    fn (ch: channel_state) -> void
    + marks the channel closed; pending receivers drain then see none
    # lifecycle
  pipeline.map_stage
    fn (in_ch: channel_state, f: fn(i64) -> i64) -> channel_state
    + returns a new channel that emits f(item) for each input item
    + closes the output channel when the input channel drains
    # stage
  pipeline.filter_stage
    fn (in_ch: channel_state, pred: fn(i64) -> bool) -> channel_state
    + returns a new channel that forwards only items satisfying pred
    # stage
  pipeline.merge
    fn (channels: list[channel_state]) -> channel_state
    + fans in multiple channels into one
    + closes the output when every input has closed
    # stage
  pipeline.fan_out
    fn (in_ch: channel_state, n: i32) -> list[channel_state]
    + distributes items round-robin across n output channels
    # stage
