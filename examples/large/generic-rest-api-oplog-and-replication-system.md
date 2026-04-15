# Requirement: "a generic oplog/replication system for REST APIs"

An append-only operation log that clients can tail to replicate state changes. Transport is pluggable; storage is pluggable.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.encoding
    std.encoding.hex_encode
      fn (data: bytes) -> string
      + encodes bytes as lowercase hex
      # encoding
  std.hash
    std.hash.sha1
      fn (data: bytes) -> bytes
      + returns 20-byte SHA-1 digest
      # hashing
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + serializes a string-to-string map as JSON
      # serialization
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization

oplog
  oplog.new
    fn (capacity: i64) -> oplog_state
    + creates an oplog with the given retention capacity
    ? capacity is the maximum number of retained ops; older ops are pruned
    # construction
  oplog.append
    fn (state: oplog_state, event_type: string, object_id: string, payload: map[string, string]) -> tuple[string, oplog_state]
    + appends an op and returns its generated id and updated state
    + the id encodes a monotonic timestamp so ids sort by insertion order
    # append
    -> std.time.now_millis
    -> std.hash.sha1
    -> std.encoding.hex_encode
  oplog.since
    fn (state: oplog_state, cursor: string) -> result[list[string], string]
    + returns ids of ops strictly after the given cursor id in order
    + returns all ids when cursor is empty
    - returns error when cursor references an id already pruned from retention
    # tail
  oplog.get
    fn (state: oplog_state, id: string) -> result[map[string, string], string]
    + returns the stored op envelope (type, object_id, payload, ts) for an id
    - returns error when the id is unknown or pruned
    # read
  oplog.prune_before
    fn (state: oplog_state, cursor: string) -> oplog_state
    + drops all ops strictly older than the given cursor
    ? no-op when cursor is empty or older than the oldest retained op
    # retention
  oplog.encode_envelope
    fn (event_type: string, object_id: string, payload: map[string, string], id: string, ts_millis: i64) -> string
    + serializes an op envelope as JSON for transport over HTTP
    # wire_format
    -> std.json.encode_object
  oplog.decode_envelope
    fn (raw: string) -> result[map[string, string], string]
    + parses a wire envelope back into its fields
    - returns error on malformed JSON or missing required fields
    # wire_format
    -> std.json.parse_object
  oplog.apply_remote
    fn (state: oplog_state, envelope: map[string, string]) -> result[oplog_state, string]
    + applies a remote op to local state if its id is newer than the last seen
    ? duplicates are ignored idempotently
    - returns error when the envelope is missing an id or ts
    # replication
  oplog.last_cursor
    fn (state: oplog_state) -> string
    + returns the id of the most recent op, or "" when empty
    # cursor
