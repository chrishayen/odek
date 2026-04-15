# Requirement: "distributed sync using operational transformation"

Clients edit a shared document concurrently. Each client sends ops; a server serializes them into a canonical stream, transforming each incoming op against any ops it missed.

std
  std.net
    std.net.ws_listen
      fn (addr: string, handler: ws_handler) -> result[server_handle, string]
      + accepts websocket connections
      # networking
    std.net.ws_connect
      fn (url: string) -> result[ws_conn, string]
      + opens a websocket
      # networking
    std.net.ws_send
      fn (conn: ws_conn, msg: string) -> result[void, string]
      + sends a text frame
      # networking
    std.net.ws_recv
      fn (conn: ws_conn) -> result[string, string]
      + receives the next text frame
      # networking
  std.json
    std.json.encode
      fn (value: json_value) -> string
      + serializes json
      # serialization
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses json
      # serialization

ot_sync
  ot_sync.make_insert
    fn (position: i32, text: string) -> op
    + returns an insert op at the given position
    # ops
  ot_sync.make_delete
    fn (position: i32, length: i32) -> op
    + returns a delete op covering length characters starting at position
    # ops
  ot_sync.apply
    fn (doc: string, op: op) -> result[string, string]
    + returns the document after applying op
    - returns error when op's range is out of bounds
    # apply
  ot_sync.transform
    fn (a: op, b: op) -> tuple[op, op]
    + returns (a', b') such that apply(apply(doc, a), b') == apply(apply(doc, b), a')
    + preserves intent under concurrent insert/insert, insert/delete, delete/delete
    ? tie-breaking for identical-position inserts uses a deterministic client id ordering
    # transform
  ot_sync.transform_against_sequence
    fn (incoming: op, history: list[op]) -> op
    + transforms incoming against every op in history in order
    # transform
  ot_sync.server_new
    fn () -> server_doc_state
    + returns an empty shared document
    # server
  ot_sync.server_receive
    fn (state: server_doc_state, client_op: op, client_revision: i64) -> result[tuple[op, server_doc_state], string]
    + transforms client_op past all ops applied since client_revision, appends it, returns the transformed op
    - returns error when client_revision is newer than the server revision
    # server
  ot_sync.server_broadcast
    fn (state: server_doc_state, outgoing: op, revision: i64) -> result[void, string]
    + pushes (outgoing, revision) to every connected client
    # server
    -> std.net.ws_send
    -> std.json.encode
  ot_sync.server_serve
    fn (state: server_doc_state, addr: string) -> result[server_handle, string]
    + listens for client websocket connections and dispatches incoming ops
    # server
    -> std.net.ws_listen
    -> std.json.parse
  ot_sync.client_new
    fn () -> client_doc_state
    + returns an empty client document
    # client
  ot_sync.client_local_edit
    fn (state: client_doc_state, op: op) -> result[client_doc_state, string]
    + applies a local op and queues it for send
    # client
  ot_sync.client_receive_remote
    fn (state: client_doc_state, remote_op: op, revision: i64) -> result[client_doc_state, string]
    + transforms pending local ops against the remote op, applies it, and advances the revision
    # client
  ot_sync.client_connect
    fn (state: client_doc_state, url: string) -> result[client_doc_state, string]
    + opens a websocket to the server and begins send/receive loop
    # client
    -> std.net.ws_connect
    -> std.net.ws_recv
