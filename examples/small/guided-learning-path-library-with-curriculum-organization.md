# Requirement: "a guided learning path library that organizes free and premium resources into an ordered curriculum"

std
  std.io
    std.io.read_file
      @ (path: string) -> result[string, string]
      + returns the full contents of a file
      - returns error when the file cannot be read
      # filesystem

learning_path
  learning_path.load
    @ (path: string) -> result[path_state, string]
    + parses a curriculum definition from a file into ordered lessons
    - returns error on malformed input
    # loading
    -> std.io.read_file
  learning_path.lessons
    @ (state: path_state) -> list[lesson]
    + returns lessons in curriculum order, each tagged free or premium
    # query
  learning_path.next_lesson
    @ (state: path_state, completed: list[string]) -> optional[lesson]
    + returns the next uncompleted lesson in order
    - returns none when every lesson has been completed
    # progression
  learning_path.filter_free
    @ (state: path_state) -> list[lesson]
    + returns only the lessons marked free
    # filtering
