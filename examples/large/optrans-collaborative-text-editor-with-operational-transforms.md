# Requirement: "a collaborative text editing service using operational transforms"

A multi-user editing backend: clients submit local edit operations against a known revision, the server transforms them against concurrent edits and broadcasts the result. std supplies monotonic time and ID generation.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.new_uuid
      fn () -> string
      + returns a random 128-bit UUID in canonical form
      # identifiers

optrans
  optrans.new_document
    fn (initial_text: string) -> document_state
    + returns a document at revision 0 containing initial_text
    # construction
  optrans.op_retain
    fn (count: i64) -> op
    + returns a retain op over the given number of characters
    # op_construction
  optrans.op_insert
    fn (text: string) -> op
    + returns an insert op for the given text
    # op_construction
  optrans.op_delete
    fn (count: i64) -> op
    + returns a delete op removing the given number of characters
    # op_construction
  optrans.apply
    fn (doc: document_state, ops: list[op]) -> result[document_state, string]
    + returns a new document with the ops applied and the revision incremented
    - returns error when total retain+delete length does not match the current document length
    # editing
  optrans.transform
    fn (a: list[op], b: list[op]) -> result[tuple[list[op], list[op]], string]
    + returns (a_prime, b_prime) such that applying a then b_prime equals applying b then a_prime
    - returns error when the two op sequences are not over the same base length
    # transformation
  optrans.compose
    fn (first: list[op], second: list[op]) -> result[list[op], string]
    + returns an op sequence equivalent to applying first then second
    - returns error when second cannot be composed onto first's output length
    # transformation
  optrans.new_session
    fn (doc: document_state) -> session_state
    + returns a session wrapping the given document with no clients
    # session_management
    -> std.id.new_uuid
  optrans.join_client
    fn (session: session_state, display_name: string) -> tuple[string, session_state]
    + returns (client_id, new_session) with the client registered at the current revision
    # session_management
    -> std.id.new_uuid
    -> std.time.now_millis
  optrans.leave_client
    fn (session: session_state, client_id: string) -> session_state
    + returns a session with the given client removed
    + leaves the session unchanged when the client is not present
    # session_management
  optrans.submit_edit
    fn (session: session_state, client_id: string, base_revision: i64, ops: list[op]) -> result[tuple[list[op], session_state], string]
    + returns (transformed_ops, new_session) after rebasing ops onto the current document
    - returns error when base_revision is greater than the current revision
    - returns error when the client is not a member of the session
    # editing
  optrans.current_text
    fn (session: session_state) -> string
    + returns the current document text
    # introspection
  optrans.current_revision
    fn (session: session_state) -> i64
    + returns the current document revision
    # introspection
