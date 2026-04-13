# Requirement: "an embedded persistent document database with a query API"

Documents are JSON-like objects. Each mutation is appended to a log; the in-memory state is rebuilt on load.

std
  std.fs
    std.fs.open_append
      @ (path: string) -> result[file_handle, string]
      + opens a file for append, creating it when missing
      # filesystem
    std.fs.write_all
      @ (handle: file_handle, data: bytes) -> result[void, string]
      + writes and flushes the buffer
      # filesystem
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the entire contents of a file
      # filesystem
    std.fs.rename
      @ (from: string, to: string) -> result[void, string]
      + atomically renames a file
      # filesystem
  std.json
    std.json.encode_value
      @ (value: json_value) -> string
      + serializes a JSON value to text
      # serialization
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses JSON text into a value
      - returns error on malformed input
      # serialization

docdb
  docdb.open
    @ (path: string) -> result[doc_state, string]
    + opens the database, replaying the append-only log into memory
    - returns error when the file is corrupt
    # lifecycle
    -> std.fs.read_all
    -> std.json.parse_value
  docdb.insert
    @ (state: doc_state, doc: json_value) -> result[tuple[string, doc_state], string]
    + assigns an id, appends the insertion record, and returns the new id
    # mutation
    -> std.fs.open_append
    -> std.fs.write_all
    -> std.json.encode_value
  docdb.update
    @ (state: doc_state, id: string, doc: json_value) -> result[doc_state, string]
    + replaces the document with the given id
    - returns error when the id is unknown
    # mutation
    -> std.fs.write_all
    -> std.json.encode_value
  docdb.remove
    @ (state: doc_state, id: string) -> result[doc_state, string]
    + appends a tombstone and drops the document from memory
    - returns error when the id is unknown
    # mutation
    -> std.fs.write_all
  docdb.find_by_id
    @ (state: doc_state, id: string) -> optional[json_value]
    + returns the document with the given id
    # query
  docdb.find_where
    @ (state: doc_state, field: string, equals: json_value) -> list[json_value]
    + returns documents whose named field equals the given value
    # query
  docdb.compact
    @ (state: doc_state, path: string) -> result[doc_state, string]
    + rewrites the log with only live documents
    # storage
    -> std.fs.write_all
    -> std.fs.rename
