# Requirement: "a puzzle challenge library with graded attempts"

Input is a marketing line, not a library. Best-effort interpretation: a puzzle store where callers submit answers and receive a grade.

std: (all units exist)

puzzle_lib
  puzzle_lib.new
    fn () -> puzzle_store
    + creates an empty store with no puzzles
    # construction
  puzzle_lib.add_puzzle
    fn (store: puzzle_store, id: string, prompt: string, expected: string) -> result[void, string]
    + registers a puzzle with a unique id, a prompt, and the expected answer
    - returns error when the id is already taken
    # registration
  puzzle_lib.get_prompt
    fn (store: puzzle_store, id: string) -> optional[string]
    + returns the prompt for the puzzle, or none if unknown
    # retrieval
  puzzle_lib.submit
    fn (store: puzzle_store, id: string, answer: string) -> result[bool, string]
    + compares the answer to the expected value after trimming and case-folding
    + returns true when the answer matches
    - returns error when the puzzle id is unknown
    # grading
  puzzle_lib.hint
    fn (store: puzzle_store, id: string, index: i32) -> optional[string]
    + returns the Nth registered hint for the puzzle, if present
    # assistance
