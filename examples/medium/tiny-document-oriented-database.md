# Requirement: "a tiny document-oriented database"

An in-memory store of JSON-like documents organized into tables, with queries by predicate and durable save/load.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads a file's full contents
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes contents, replacing any existing file
      # filesystem
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON string into a dynamic value
      - returns error on malformed JSON
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + serializes a dynamic value to a JSON string
      # serialization

document_db
  document_db.new
    @ () -> db_state
    + creates an empty database with no tables
    # construction
  document_db.table
    @ (state: db_state, name: string) -> db_state
    + ensures a table with the given name exists
    # tables
  document_db.insert
    @ (state: db_state, table: string, doc: json_value) -> result[tuple[db_state, i64], string]
    + inserts the document into the table and returns a newly assigned id
    - returns error when the table does not exist
    # insert
  document_db.get
    @ (state: db_state, table: string, id: i64) -> result[optional[json_value], string]
    + returns the document with the given id, or none if removed
    - returns error when the table does not exist
    # get
  document_db.update
    @ (state: db_state, table: string, id: i64, doc: json_value) -> result[db_state, string]
    + replaces the document at the given id
    - returns error when the id does not exist in the table
    # update
  document_db.remove
    @ (state: db_state, table: string, id: i64) -> result[db_state, string]
    + deletes the document with the given id
    - returns error when the id does not exist in the table
    # remove
  document_db.search
    @ (state: db_state, table: string, predicate: doc_predicate) -> result[list[json_value], string]
    + returns all documents in the table for which the predicate returns true
    - returns error when the table does not exist
    # search
  document_db.save
    @ (state: db_state, path: string) -> result[void, string]
    + writes the full database to disk as a single JSON document
    # persistence
    -> std.json.encode
    -> std.fs.write_all
  document_db.load
    @ (path: string) -> result[db_state, string]
    + reads a database from disk
    - returns error when the file is missing or malformed
    # persistence
    -> std.fs.read_all
    -> std.json.parse
