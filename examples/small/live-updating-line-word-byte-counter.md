# Requirement: "a live-updating line/word/byte counter"

A counter that can be updated incrementally as new bytes arrive and yields running totals.

std: (all units exist)

wc
  wc.new
    fn () -> wc_state
    + creates an empty counter with zero lines, words, and bytes
    # construction
  wc.feed
    fn (state: wc_state, chunk: bytes) -> wc_state
    + updates counts for newlines, whitespace-delimited words, and total bytes
    + correctly handles words split across chunk boundaries
    - does not double-count a word when a chunk ends mid-word
    # incremental_counting
  wc.totals
    fn (state: wc_state) -> tuple[i64, i64, i64]
    + returns (lines, words, bytes) as of the most recent feed
    # snapshot
