# Requirement: "a library that filters duplicate lines from unsorted input"

Streams lines and emits each only on its first occurrence.

std: (all units exist)

uniq
  uniq.new
    fn () -> uniq_state
    + creates an empty deduplication state
    # construction
  uniq.observe
    fn (state: uniq_state, line: string) -> tuple[bool, uniq_state]
    + returns (true, new_state) when line has not been seen before
    - returns (false, unchanged_state) when line has already been seen
    # deduplication
  uniq.filter_lines
    fn (lines: list[string]) -> list[string]
    + returns lines in original order with later duplicates removed
    # batch_api
