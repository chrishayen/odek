# Requirement: "a federated blogging backend that speaks the ActivityPub protocol"

Local authoring plus the subset of the federation protocol needed to publish posts and receive follows. Storage and HTTP transport are injected.

std
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a flat JSON object into a string map
      - returns error on malformed JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a flat JSON object
      # serialization
  std.http
    std.http.post_json
      fn (url: string, headers: map[string, string], body: string) -> result[http_response, string]
      + sends a POST with the given JSON body
      - returns error on network failure
      # http
    std.http.get
      fn (url: string, headers: map[string, string]) -> result[http_response, string]
      + sends a GET
      - returns error on network failure
      # http
  std.crypto
    std.crypto.rsa_generate
      fn (bits: i32) -> result[key_pair, string]
      + generates an RSA key pair
      # cryptography
    std.crypto.sign_rsa_sha256
      fn (key: key_pair, data: bytes) -> bytes
      + returns an RSA-SHA256 signature of data
      # cryptography
    std.crypto.verify_rsa_sha256
      fn (public_key: bytes, data: bytes, signature: bytes) -> bool
      + returns true when the signature matches
      # cryptography
  std.time
    std.time.now_rfc3339
      fn () -> string
      + returns the current instant in RFC-3339 format
      # time

blog
  blog.actor_new
    fn (username: string, display_name: string, host: string) -> result[actor, string]
    + creates a local actor with a fresh key pair and a canonical id URL
    # identity
    -> std.crypto.rsa_generate
  blog.actor_document
    fn (a: actor) -> string
    + returns the JSON-LD document describing the actor
    # identity
    -> std.json.encode_object
  blog.post_create
    fn (author: actor, title: string, body: string) -> post
    + creates an unpublished post owned by the actor
    # authoring
    -> std.time.now_rfc3339
  blog.post_to_activity
    fn (p: post) -> string
    + serializes a post as a "Create" activity with a nested "Note" object
    # activity
    -> std.json.encode_object
    -> std.time.now_rfc3339
  blog.inbox_receive
    fn (recipient: actor, raw_body: string, signature_header: string) -> result[activity, string]
    + parses an incoming activity and verifies its HTTP signature
    - returns error when the signature does not match the sender's public key
    - returns error when the body is not valid JSON
    # federation
    -> std.json.parse_object
    -> std.crypto.verify_rsa_sha256
  blog.outbox_deliver
    fn (sender: actor, recipient_inbox: string, activity_json: string) -> result[void, string]
    + signs and POSTs the activity to the recipient's inbox
    - returns error when the HTTP call fails
    # federation
    -> std.crypto.sign_rsa_sha256
    -> std.http.post_json
  blog.follow_accept
    fn (local: actor, follow_activity: activity) -> result[activity, string]
    + produces an "Accept" activity acknowledging a follow
    # social
    -> std.json.encode_object
  blog.webfinger_resolve
    fn (handle: string) -> result[string, string]
    + resolves "user@host" to an actor URL via WebFinger
    - returns error when the handle is malformed
    - returns error when the remote host does not respond
    # discovery
    -> std.http.get
    -> std.json.parse_object
  blog.fetch_remote_actor
    fn (actor_url: string) -> result[remote_actor, string]
    + fetches and parses a remote actor document
    - returns error when the document cannot be parsed
    # discovery
    -> std.http.get
    -> std.json.parse_object
  blog.timeline_append
    fn (local: actor, activity: activity) -> result[void, string]
    + stores an inbound activity in the recipient's timeline via the injected storage driver
    # storage
