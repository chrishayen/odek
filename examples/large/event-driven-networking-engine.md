# Requirement: "an event-driven networking engine"

A reactor-style loop managing sockets, timers, and user callbacks. Transport primitives live in std; the project layer is the scheduler and dispatch.

std
  std.net
    std.net.listen_tcp
      fn (host: string, port: i32) -> result[listener, string]
      + binds and starts listening on the given address
      - returns error when the port cannot be bound
      # network
    std.net.accept
      fn (l: listener) -> result[socket, string]
      + accepts one pending connection
      - returns error when the listener has been closed
      # network
    std.net.set_nonblocking
      fn (s: socket, nonblocking: bool) -> result[void, string]
      + toggles nonblocking mode on a socket
      # network
    std.net.read_available
      fn (s: socket) -> result[bytes, string]
      + reads bytes that are ready on the socket without blocking
      - returns empty bytes when nothing is buffered
      - returns error when the peer has closed
      # network
    std.net.write_bytes
      fn (s: socket, data: bytes) -> result[i64, string]
      + writes up to the available send buffer and returns bytes written
      - returns error on broken connection
      # network
    std.net.poll_ready
      fn (sockets: list[socket], timeout_millis: i64) -> list[ready_event]
      + blocks until a socket is readable or writable or timeout elapses
      # network
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time

netloop
  netloop.new
    fn () -> loop_state
    + creates an empty reactor with no registrations
    # construction
  netloop.register_listener
    fn (state: loop_state, l: listener, accept_tag: string) -> loop_state
    + registers a listener whose accepted sockets will be handed to the named accept callback
    # registration
  netloop.register_socket
    fn (state: loop_state, s: socket, read_tag: string, write_tag: string) -> loop_state
    + registers a socket with callbacks for read and write events
    -> std.net.set_nonblocking
    # registration
  netloop.schedule_timer
    fn (state: loop_state, delay_millis: i64, callback_tag: string) -> loop_state
    + enqueues a timer that fires the named callback after the given delay
    -> std.time.now_millis
    # timing
  netloop.cancel_timer
    fn (state: loop_state, timer_id: i64) -> loop_state
    + removes a pending timer by id
    # timing
  netloop.unregister_socket
    fn (state: loop_state, s: socket) -> loop_state
    + removes a socket from the reactor
    # registration
  netloop.tick
    fn (state: loop_state) -> result[tuple[loop_state, list[dispatch], i64], string]
    + runs one poll iteration and returns the new state, pending dispatches, and the next timer deadline
    -> std.net.poll_ready
    -> std.time.now_millis
    # dispatch
  netloop.ready_event
    fn (kind: string, socket_id: i64) -> ready_event
    + builds a ready event descriptor
    # events
  netloop.buffered_send
    fn (state: loop_state, s: socket, data: bytes) -> loop_state
    + appends bytes to a socket's outgoing buffer, to be drained by the reactor
    # buffering
  netloop.drain_writes
    fn (state: loop_state, s: socket) -> result[loop_state, string]
    + attempts to flush a socket's outgoing buffer
    -> std.net.write_bytes
    # buffering
  netloop.drain_reads
    fn (state: loop_state, s: socket) -> result[tuple[loop_state, bytes], string]
    + reads available bytes into an inbound buffer and returns them
    -> std.net.read_available
    # buffering
