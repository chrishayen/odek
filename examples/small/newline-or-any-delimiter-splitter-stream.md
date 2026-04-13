# Requirement: "a streaming delimiter splitter"

Incremental splitting: feed bytes in, get complete records out, keep the tail for the next call.

std: (all units exist)

splitter
  splitter.new
    @ (delimiter: bytes) -> splitter_state
    + creates a splitter for the given delimiter
    - returns a splitter with empty delimiter that never splits when delimiter is empty
    # construction
  splitter.feed
    @ (state: splitter_state, chunk: bytes) -> tuple[list[bytes], splitter_state]
    + returns complete records found in chunk plus buffered tail
    + handles delimiters that straddle chunk boundaries
    - returns empty list when no delimiter is present yet
    # streaming
  splitter.flush
    @ (state: splitter_state) -> bytes
    + returns any remaining buffered bytes without a trailing delimiter
    # termination
