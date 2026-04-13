# Requirement: "a multiplayer code editor"

A collaborative text buffer with syntax highlighting, cursor tracking per participant, and CRDT-based concurrent edits. This library is the editor core, not a UI toolkit or networking layer — it produces state transitions the caller serializes and transports.

std
  std.time
    std.time.now_nanos
      @ () -> i64
      + returns current unix time in nanoseconds
      # time
  std.uuid
    std.uuid.new_v4
      @ () -> string
      + returns a random UUID as a string
      # identifiers
  std.hash
    std.hash.hash64
      @ (data: bytes) -> u64
      + returns a 64-bit non-cryptographic hash
      # hashing

editor
  editor.new_document
    @ (initial_text: string) -> document_state
    + creates an empty document or seeds with initial text
    # construction
  editor.join_participant
    @ (state: document_state, display_name: string) -> tuple[string, document_state]
    + assigns a participant id and records their presence
    # presence
    -> std.uuid.new_v4
  editor.leave_participant
    @ (state: document_state, participant_id: string) -> document_state
    + removes the participant and their cursors
    # presence
  editor.local_insert
    @ (state: document_state, participant_id: string, offset: i64, text: string) -> tuple[edit_op, document_state]
    + applies an insert locally and returns the op to broadcast
    - returns unchanged state when participant is unknown
    # editing
    -> std.time.now_nanos
    -> std.uuid.new_v4
  editor.local_delete
    @ (state: document_state, participant_id: string, offset: i64, length: i64) -> tuple[edit_op, document_state]
    + applies a delete locally and returns the op to broadcast
    - returns unchanged state when participant is unknown
    # editing
    -> std.time.now_nanos
  editor.apply_remote_op
    @ (state: document_state, op: edit_op) -> document_state
    + integrates a remote op using CRDT positions
    ? idempotent: re-applying the same op is a no-op
    # crdt
    -> std.hash.hash64
  editor.move_cursor
    @ (state: document_state, participant_id: string, offset: i64) -> document_state
    + updates a participant's cursor position
    # presence
  editor.cursors
    @ (state: document_state) -> list[cursor_view]
    + returns all current cursor positions with participant ids
    # presence
  editor.get_text
    @ (state: document_state) -> string
    + returns the current document text
    # read
  editor.new_language
    @ (name: string, keywords: list[string], comment_prefix: string) -> language_spec
    + defines a language with keywords and line-comment syntax
    # language
  editor.highlight
    @ (state: document_state, lang: language_spec) -> list[highlight_token]
    + returns token spans classified as keyword, comment, identifier, or literal
    # syntax_highlighting
  editor.undo
    @ (state: document_state, participant_id: string) -> document_state
    + reverses the participant's most recent local op
    - no-op when participant has no undo history
    # history
  editor.redo
    @ (state: document_state, participant_id: string) -> document_state
    + re-applies the participant's most recent undone op
    - no-op when redo history is empty
    # history
  editor.diff_since
    @ (state: document_state, version: i64) -> list[edit_op]
    + returns ops applied since the given version for catch-up
    # synchronization
