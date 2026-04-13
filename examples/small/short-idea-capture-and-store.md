# Requirement: "a library to capture and store short ideas"

Append-only idea journal with a title, body, and timestamp per entry.

std
  std.fs
    std.fs.append
      @ (path: string, content: string) -> result[void, string]
      + appends content to the file, creating it if needed
      - returns error on write failure
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns full file contents
      - returns error when unreadable
      # filesystem
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

ideas
  ideas.capture
    @ (store_path: string, title: string, body: string) -> result[void, string]
    + appends one record with the given title, body, and current timestamp
    - returns error when title is empty
    # capture
    -> std.time.now_seconds
    -> std.fs.append
  ideas.list_all
    @ (store_path: string) -> result[list[idea_entry], string]
    + returns every stored idea in chronological order
    + returns an empty list when the store is missing
    # listing
    -> std.fs.read_all
  ideas.search
    @ (store_path: string, query: string) -> result[list[idea_entry], string]
    + returns entries whose title or body contains query, case-insensitive
    - returns error when the store is corrupt
    # search
    -> std.fs.read_all
