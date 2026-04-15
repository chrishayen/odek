# Requirement: "a todo list library with file persistence"

Six project operations. `save` and `load` wire to std filesystem and JSON primitives that any data-at-rest project would reuse.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file at path into bytes
      - returns error when the file does not exist
      - returns error when the file is not readable
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data to path, creating the file if missing, truncating otherwise
      - returns error when the path is not writable
      # filesystem
  std.json
    std.json.encode
      fn (value: any) -> string
      + serializes any supported type to JSON
      # serialization
    std.json.decode
      fn (raw: string, type: type) -> result[any, string]
      + parses JSON into the target type
      - returns error on malformed JSON
      - returns error when the JSON does not match the target type
      # serialization

todo
  todo.add
    fn (state: todo_state, item: string) -> todo_state
    + appends a new item with a fresh integer id
    ? ids are monotonically increasing; deleted ids are not reused
    # state_mutation
  todo.complete
    fn (state: todo_state, id: i32) -> todo_state
    + marks the item with the given id as completed
    ? completing a non-existent id is a no-op
    # state_mutation
  todo.remove
    fn (state: todo_state, id: i32) -> todo_state
    + removes the item with the given id
    ? removing a non-existent id is a no-op
    # state_mutation
  todo.list
    fn (state: todo_state) -> list[todo_item]
    + returns all items in insertion order
    # state_access
  todo.save
    fn (state: todo_state, path: string) -> result[void, string]
    + serializes state to JSON and writes it to the given path
    - returns error when the file cannot be written
    # persistence
    -> std.json.encode
    -> std.fs.write_all
  todo.load
    fn (path: string) -> result[todo_state, string]
    + reads the file at path and parses it as todo state
    + returns an empty state when the file does not exist
    - returns error when the contents are not valid JSON
    # persistence
    -> std.fs.read_all
    -> std.json.decode
