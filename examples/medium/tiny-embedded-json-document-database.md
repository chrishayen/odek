# Requirement: "a tiny embedded json document database"

In-memory key-value document store with load and save hooks. Persistence uses std filesystem primitives so it can be swapped in tests.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns the file contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, content: string) -> result[void, string]
      + writes content, creating or replacing the file
      - returns error on permission failure
      # filesystem
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

doc_db
  doc_db.new
    @ () -> doc_db_state
    + returns an empty document store
    # construction
  doc_db.set
    @ (db: doc_db_state, key: string, value: string) -> doc_db_state
    + returns a store with the value written at the key
    # write
  doc_db.get
    @ (db: doc_db_state, key: string) -> optional[string]
    + returns the stored value
    - returns none when the key is absent
    # read
  doc_db.remove
    @ (db: doc_db_state, key: string) -> doc_db_state
    + returns a store with the key removed
    # write
  doc_db.load
    @ (path: string) -> result[doc_db_state, string]
    + reads and parses the store from disk
    - returns error when the file is missing or invalid
    # persistence
    -> std.fs.read_all
    -> std.json.parse_object
  doc_db.save
    @ (db: doc_db_state, path: string) -> result[void, string]
    + serializes and writes the store to disk
    - returns error on write failure
    # persistence
    -> std.json.encode_object
    -> std.fs.write_all
