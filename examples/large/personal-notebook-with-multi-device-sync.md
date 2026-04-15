# Requirement: "a personal notebook with multi-device sync"

Notes live in local books; a sync protocol exchanges an append-only action log with a remote endpoint, merging by timestamp.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire file into memory
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes data atomically replacing any existing file
      # filesystem
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns entries in the directory, not including dot entries
      # filesystem
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document into a dynamic value
      # serialization
    std.json.encode
      fn (value: json_value) -> string
      + encodes a dynamic value as a JSON string
      # serialization
  std.http
    std.http.post_json
      fn (url: string, headers: map[string, string], body: string) -> result[http_response, string]
      + performs an HTTPS POST with the given headers and JSON body
      - returns error on network failure
      # http
  std.time
    std.time.now_millis
      fn () -> i64
      + returns the current time in milliseconds since the epoch
      # time
  std.uuid
    std.uuid.new_v4
      fn () -> string
      + returns a random UUID in canonical form
      # identity

notebook
  notebook.open
    fn (data_dir: string) -> result[notebook_state, string]
    + loads all books and pending actions from the data directory
    - returns error when the directory cannot be read
    # lifecycle
    -> std.fs.list_dir
    -> std.fs.read_all
    -> std.json.parse
  notebook.create_book
    fn (state: notebook_state, name: string) -> result[notebook_state, string]
    + adds a new empty book
    - returns error when a book with that name already exists
    # books
  notebook.add_note
    fn (state: notebook_state, book: string, body: string) -> result[notebook_state, string]
    + appends a note to the given book and records an add action
    - returns error when the book does not exist
    # notes
    -> std.uuid.new_v4
    -> std.time.now_millis
  notebook.edit_note
    fn (state: notebook_state, note_id: string, body: string) -> result[notebook_state, string]
    + replaces a note's body and records an edit action
    - returns error when the note id is unknown
    # notes
    -> std.time.now_millis
  notebook.remove_note
    fn (state: notebook_state, note_id: string) -> result[notebook_state, string]
    + marks the note as removed and records a delete action
    # notes
    -> std.time.now_millis
  notebook.list_notes
    fn (state: notebook_state, book: string) -> list[note]
    + returns non-deleted notes in the given book in creation order
    # notes
  notebook.save
    fn (state: notebook_state) -> result[void, string]
    + writes books and the pending action log back to the data directory
    # persistence
    -> std.json.encode
    -> std.fs.write_all
  notebook.sync
    fn (state: notebook_state, remote_url: string, token: string) -> result[notebook_state, string]
    + uploads pending actions, downloads remote actions, and merges them by timestamp
    - returns error on transport failure, leaving local state unchanged
    # sync
    -> std.http.post_json
    -> std.json.encode
    -> std.json.parse
  notebook.merge_actions
    fn (state: notebook_state, incoming: list[sync_action]) -> notebook_state
    + applies incoming actions that are newer than the local state for each affected note
    + tolerates duplicates by ignoring already-seen action ids
    # sync
  notebook.diff_since
    fn (state: notebook_state, since_millis: i64) -> list[sync_action]
    + returns actions recorded at or after the given timestamp
    # sync
