# Requirement: "a course lesson index"

The input here is a video series for learning programming and game development. The minimal library-shaped interpretation is an ordered lesson index keyed by lesson number.

std: (all units exist)

lesson_index
  lesson_index.new
    fn () -> lesson_index_state
    + returns an empty lesson index
    # construction
  lesson_index.add_lesson
    fn (state: lesson_index_state, number: i32, title: string) -> lesson_index_state
    + appends a lesson at the given number, preserving numeric order
    + replaces the title when number already exists
    # writes
  lesson_index.lookup
    fn (state: lesson_index_state, number: i32) -> optional[string]
    + returns the title for the given lesson number
    - returns none when lesson number does not exist
    # reads
