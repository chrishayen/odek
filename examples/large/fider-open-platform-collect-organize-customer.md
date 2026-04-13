# Requirement: "a customer feedback collection and organization library"

Users submit ideas, vote on them, and comment; administrators change status and tag ideas. Persistence is delegated to a std key-value store.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.kv
    std.kv.put
      @ (key: string, value: bytes) -> result[void, string]
      + writes a key-value pair
      # storage
    std.kv.get
      @ (key: string) -> result[optional[bytes], string]
      + returns the stored value or none when absent
      # storage
    std.kv.list_prefix
      @ (prefix: string) -> result[list[string], string]
      + returns every key starting with prefix
      # storage
  std.id
    std.id.random_uuid
      @ () -> string
      + returns a fresh UUID v4
      # identifiers
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + serializes a string-to-string map
      # serialization
    std.json.parse_object
      @ (raw: string) -> result[map[string, string], string]
      + parses a JSON object
      - returns error on invalid JSON
      # serialization

feedback
  feedback.register_user
    @ (name: string, email: string) -> result[user_id, string]
    + creates a user record and returns its id
    - returns error when email is already registered
    # users
    -> std.id.random_uuid
    -> std.kv.put
  feedback.create_idea
    @ (author: user_id, title: string, body: string) -> result[idea_id, string]
    + creates an idea in status "open" with a fresh id
    - returns error when title is empty
    # ideas
    -> std.id.random_uuid
    -> std.time.now_seconds
    -> std.json.encode_object
    -> std.kv.put
  feedback.vote
    @ (voter: user_id, idea: idea_id) -> result[i32, string]
    + records a vote and returns the new total
    - returns error when voter has already voted on the idea
    # ideas
    -> std.kv.get
    -> std.kv.put
  feedback.comment
    @ (author: user_id, idea: idea_id, body: string) -> result[comment_id, string]
    + appends a comment to the idea
    - returns error when the idea does not exist
    # comments
    -> std.id.random_uuid
    -> std.time.now_seconds
    -> std.kv.put
  feedback.change_status
    @ (idea: idea_id, status: idea_status) -> result[void, string]
    + updates the idea status (open, planned, started, completed, declined)
    - returns error when the idea does not exist
    # ideas
    -> std.kv.get
    -> std.kv.put
  feedback.tag_idea
    @ (idea: idea_id, tag: string) -> result[void, string]
    + attaches a tag to the idea; tags are a set
    # organization
    -> std.kv.get
    -> std.kv.put
  feedback.list_by_status
    @ (status: idea_status) -> result[list[idea], string]
    + returns every idea with the given status, newest first
    # queries
    -> std.kv.list_prefix
    -> std.kv.get
    -> std.json.parse_object
  feedback.search
    @ (query: string) -> result[list[idea], string]
    + returns ideas whose title or body contains the query, case-insensitive
    # queries
    -> std.kv.list_prefix
  feedback.top_voted
    @ (limit: i32) -> result[list[idea], string]
    + returns the top N open ideas by vote count
    # queries
