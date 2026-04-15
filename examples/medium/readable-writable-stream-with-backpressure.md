# Requirement: "a readable/writable stream abstraction with backpressure"

A stream holds an internal buffer. Writers push chunks; readers drain them. The stream reports when it is full and when it has ended.

std: (all units exist)

streams
  streams.new_readable
    fn (high_water_mark: i32) -> readable_state
    + returns an empty readable stream with the given buffer threshold
    # construction
  streams.push
    fn (state: readable_state, chunk: bytes) -> tuple[bool, readable_state]
    + returns (accept_more, new_state) where accept_more is false when the buffer exceeds the high water mark
    - rejects chunks after end_of_stream has been signaled, leaving state unchanged
    # backpressure
  streams.end_of_stream
    fn (state: readable_state) -> readable_state
    + marks the stream as finished; subsequent reads drain remaining buffer then return empty
    # lifecycle
  streams.read
    fn (state: readable_state, max_bytes: i32) -> tuple[bytes, readable_state]
    + returns up to max_bytes and a state with those bytes removed
    + returns empty bytes when the buffer is empty
    # draining
  streams.is_ended
    fn (state: readable_state) -> bool
    + returns true when end_of_stream was called and the buffer is empty
    # lifecycle
  streams.pipe
    fn (src: readable_state, dst: readable_state, max_bytes: i32) -> tuple[readable_state, readable_state]
    + reads up to max_bytes from src and pushes them into dst
    + propagates end_of_stream when src becomes ended
    # composition
