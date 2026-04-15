# Requirement: "a multiplexor for aligned tickers"

One ticker source that multiple subscribers share, with ticks aligned to wall-clock boundaries so all subscribers see the same timestamps.

std
  std.time
    std.time.now_nanos
      fn () -> i64
      + returns current unix time in nanoseconds
      # time
    std.time.sleep_nanos
      fn (duration: i64) -> void
      + blocks the caller for the given number of nanoseconds
      # time

ticker_mux
  ticker_mux.new
    fn (interval_nanos: i64) -> ticker_mux_state
    + creates a multiplexor with no subscribers that will tick every interval_nanos
    ? interval_nanos must be positive
    # construction
  ticker_mux.subscribe
    fn (state: ticker_mux_state) -> tuple[subscriber_id, ticker_mux_state]
    + returns a new subscriber id and the updated state
    # subscription
  ticker_mux.unsubscribe
    fn (state: ticker_mux_state, id: subscriber_id) -> ticker_mux_state
    + removes the subscriber; no-op when id is unknown
    # subscription
  ticker_mux.next_tick_nanos
    fn (state: ticker_mux_state) -> i64
    + returns the next wall-clock instant aligned to the interval boundary
    ? alignment uses floor(now / interval) * interval + interval
    # alignment
    -> std.time.now_nanos
  ticker_mux.wait_and_broadcast
    fn (state: ticker_mux_state) -> tuple[i64, list[subscriber_id]]
    + sleeps until the next aligned tick and returns the tick timestamp plus the list of subscribers to notify
    # dispatch
    -> std.time.sleep_nanos
