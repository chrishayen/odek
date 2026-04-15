# Requirement: "a relationship-based fine-grained authorization engine"

Stores tuples of (object, relation, subject) and answers permission checks based on a schema of relations and rewrites.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns the current unix time in milliseconds
      # time
  std.hash
    std.hash.fnv64
      fn (data: bytes) -> u64
      + computes FNV-1a 64-bit hash of data
      # hashing

relauth
  relauth.new_store
    fn () -> store_state
    + creates an empty relation tuple store
    # construction
  relauth.define_schema
    fn (state: store_state, schema: namespace_schema) -> store_state
    + registers a namespace with its relations and rewrite rules
    ? rewrite rules express computed relations like "editor implies viewer"
    # schema
  relauth.write_tuple
    fn (state: store_state, object: string, relation: string, subject: string) -> result[store_state, string]
    + records a relationship tuple in the store
    - returns error when the object's namespace is unknown
    # writes
    -> std.time.now_millis
  relauth.delete_tuple
    fn (state: store_state, object: string, relation: string, subject: string) -> result[store_state, string]
    + removes a relationship tuple
    - returns error when the tuple does not exist
    # writes
  relauth.check
    fn (state: store_state, object: string, relation: string, subject: string) -> result[bool, string]
    + returns true when the subject has the given relation on the object, applying rewrite rules
    - returns error when the relation is not defined in the namespace
    # authorization
  relauth.expand
    fn (state: store_state, object: string, relation: string) -> result[subject_tree, string]
    + returns the tree of subjects that resolve to the given object-relation
    - returns error when the relation is not defined
    # introspection
  relauth.lookup_resources
    fn (state: store_state, subject: string, relation: string, namespace: string) -> list[string]
    + returns all objects in a namespace where the subject has the relation
    # queries
  relauth.lookup_subjects
    fn (state: store_state, object: string, relation: string) -> list[string]
    + returns all subjects holding the relation on the object
    # queries
  relauth.snapshot_token
    fn (state: store_state) -> string
    + returns an opaque token representing the current store revision
    # consistency
    -> std.hash.fnv64
  relauth.check_at_snapshot
    fn (state: store_state, token: string, object: string, relation: string, subject: string) -> result[bool, string]
    + performs a check against the store version identified by the token
    - returns error when the token is unknown or expired
    # authorization
