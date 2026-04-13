# Requirement: "a library for managing tasks, boards, and notes"

An in-memory model for a personal task tracker with boards, tasks, and notes. Persistence is optional and delegates to std filesystem.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the full contents of the file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes bytes to the path
      # filesystem
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses a JSON value
      - returns error on invalid JSON
      # serialization
    std.json.encode_value
      @ (value: json_value) -> string
      + encodes a JSON value to text
      # serialization

trackerlib
  trackerlib.new
    @ () -> tracker_state
    + creates an empty tracker with a default "My Board" board
    # construction
  trackerlib.add_board
    @ (state: tracker_state, title: string) -> tuple[i32, tracker_state]
    + creates a board and returns its id
    # mutation
  trackerlib.add_task
    @ (state: tracker_state, board_id: i32, title: string) -> result[tuple[i32, tracker_state], string]
    + creates a task on the board
    - returns error when the board does not exist
    # mutation
  trackerlib.add_note
    @ (state: tracker_state, board_id: i32, body: string) -> result[tuple[i32, tracker_state], string]
    + creates a note on the board
    - returns error when the board does not exist
    # mutation
  trackerlib.check_task
    @ (state: tracker_state, task_id: i32) -> result[tracker_state, string]
    + marks the task completed
    - returns error when the task does not exist
    # mutation
  trackerlib.move_task
    @ (state: tracker_state, task_id: i32, to_board_id: i32) -> result[tracker_state, string]
    + moves the task to a different board
    - returns error when either id does not exist
    # mutation
  trackerlib.delete_item
    @ (state: tracker_state, item_id: i32) -> result[tracker_state, string]
    + removes the task or note
    - returns error when the id does not exist
    # mutation
  trackerlib.list_board
    @ (state: tracker_state, board_id: i32) -> result[board_view, string]
    + returns the tasks and notes on the board in insertion order
    - returns error when the board does not exist
    # query
  trackerlib.save
    @ (state: tracker_state, path: string) -> result[void, string]
    + serializes the tracker to a file
    # persistence
    -> std.json.encode_value
    -> std.fs.write_all
  trackerlib.load
    @ (path: string) -> result[tracker_state, string]
    + loads a tracker from a file
    - returns error when the file is missing or malformed
    # persistence
    -> std.fs.read_all
    -> std.json.parse_value
