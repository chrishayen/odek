# Requirement: "a line-by-line stream reader"

Wraps an incremental byte feed and yields whole lines as they arrive, buffering partial data across feeds.

std: (all units exist)

linestream
  linestream.new
    @ () -> stream_state
    + creates an empty line stream buffer
    # construction
  linestream.feed
    @ (state: stream_state, chunk: bytes) -> stream_state
    + appends bytes to the internal buffer
    # input
  linestream.next_line
    @ (state: stream_state) -> tuple[optional[string], stream_state]
    + returns the next complete line (without its trailing newline) and the updated state
    + recognizes both "\n" and "\r\n" terminators
    - returns none when the buffer holds no complete line
    # parsing
  linestream.flush
    @ (state: stream_state) -> tuple[optional[string], stream_state]
    + returns any remaining partial content as a final line and clears the buffer
    - returns none when the buffer is empty
    # termination
