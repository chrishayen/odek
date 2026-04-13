# Requirement: "a federated social server speaking the ActivityPub protocol"

Local accounts, posts, follow graph, inbox/outbox delivery, and signed ActivityPub federation. The project layer coordinates std primitives for storage, HTTP, cryptography, and JSON.

std
  std.http
    std.http.serve
      @ (addr: string, handler: fn(http_request) -> http_response) -> result[void, string]
      + starts an HTTP server that dispatches every request to handler
      - returns error when the address cannot be bound
      # networking
    std.http.post
      @ (url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + issues a POST request with the given headers and body
      - returns error on transport failure
      # networking
  std.json
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string map as JSON
      # serialization
  std.crypto
    std.crypto.generate_rsa_keypair
      @ (bits: i32) -> result[tuple[bytes, bytes], string]
      + returns (private_key, public_key) in PEM form
      - returns error when bits is too small
      # cryptography
    std.crypto.rsa_sign
      @ (private_key: bytes, data: bytes) -> result[bytes, string]
      + produces an RSA-SHA256 signature
      - returns error when the key is malformed
      # cryptography
    std.crypto.rsa_verify
      @ (public_key: bytes, data: bytes, signature: bytes) -> bool
      + returns true when the signature is valid
      - returns false on mismatch
      # cryptography
  std.sql
    std.sql.open
      @ (dsn: string) -> result[sql_conn, string]
      + opens a database connection
      - returns error on invalid DSN
      # database
    std.sql.exec
      @ (conn: sql_conn, statement: string, params: list[string]) -> result[i64, string]
      + runs a statement and returns the affected row count
      - returns error on SQL failure
      # database
    std.sql.query
      @ (conn: sql_conn, statement: string, params: list[string]) -> result[list[map[string, string]], string]
      + runs a query and returns matching rows
      - returns error on SQL failure
      # database
  std.time
    std.time.now_rfc3339
      @ () -> string
      + returns current time as an RFC 3339 string
      # time

fediverse
  fediverse.new
    @ (conn: sql_conn, domain: string) -> fediverse_state
    + creates a server state bound to a database and host domain
    # construction
  fediverse.create_account
    @ (state: fediverse_state, username: string) -> result[account, string]
    + creates a local account with a freshly generated keypair
    - returns error when the username already exists
    # accounts
    -> std.crypto.generate_rsa_keypair
    -> std.sql.exec
  fediverse.lookup_actor
    @ (state: fediverse_state, username: string) -> result[account, string]
    + looks up a local account by username
    - returns error when the account does not exist
    # accounts
    -> std.sql.query
  fediverse.render_actor
    @ (state: fediverse_state, acct: account) -> string
    + renders the actor object as ActivityPub JSON with inbox and public key
    # federation
    -> std.json.encode_object
  fediverse.create_post
    @ (state: fediverse_state, author: account, content: string) -> result[post, string]
    + creates a local post addressed to followers and the public collection
    - returns error when content is empty
    # posts
    -> std.sql.exec
    -> std.time.now_rfc3339
  fediverse.follow
    @ (state: fediverse_state, follower: account, target_actor_url: string) -> result[void, string]
    + sends a signed Follow activity to the target actor
    - returns error on delivery failure
    # follow_graph
  fediverse.sign_request
    @ (acct: account, method: string, url: string, body: bytes) -> map[string, string]
    + returns HTTP Signature headers using the account's private key
    # federation
    -> std.crypto.rsa_sign
  fediverse.verify_signature
    @ (state: fediverse_state, method: string, url: string, headers: map[string, string], body: bytes) -> result[string, string]
    + returns the verified sender actor URL
    - returns error when the signature cannot be verified
    # federation
    -> std.crypto.rsa_verify
  fediverse.deliver
    @ (state: fediverse_state, inbox_url: string, activity: string, signer: account) -> result[void, string]
    + POSTs a signed activity to a remote inbox
    - returns error on transport failure
    # federation
    -> std.http.post
  fediverse.handle_inbox
    @ (state: fediverse_state, req: http_request) -> http_response
    + verifies the signature and dispatches the activity to the right handler
    - returns 401 when the signature is invalid
    - returns 400 on malformed JSON
    # inbox
    -> std.json.parse_object
  fediverse.handle_webfinger
    @ (state: fediverse_state, resource: string) -> result[string, string]
    + returns the WebFinger JSON describing a local account
    - returns error when the resource is unknown
    # discovery
    -> std.json.encode_object
  fediverse.serve
    @ (state: fediverse_state, addr: string) -> result[void, string]
    + starts an HTTP server that routes inbox, outbox, actor, and webfinger paths
    - returns error when the address cannot be bound
    # serving
    -> std.http.serve
