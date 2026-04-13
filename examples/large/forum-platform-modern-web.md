# Requirement: "a forum platform library"

Users, threads, posts, votes, and a moderation log. Persistence and auth primitives live in std.

std
  std.sql
    std.sql.open
      @ (dsn: string) -> result[db_handle, string]
      + opens a database connection pool
      - returns error on unreachable database
      # database
    std.sql.exec
      @ (db: db_handle, query: string, params: list[sql_value]) -> result[i64, string]
      + executes a statement and returns rows-affected
      # database
    std.sql.query
      @ (db: db_handle, query: string, params: list[sql_value]) -> result[list[map[string, sql_value]], string]
      + runs a query and returns rows as maps
      # database
  std.crypto
    std.crypto.bcrypt_hash
      @ (password: string, cost: i32) -> result[string, string]
      + returns a bcrypt hash of the password
      # cryptography
    std.crypto.bcrypt_verify
      @ (password: string, hash: string) -> bool
      + returns true when the password matches the hash
      # cryptography
  std.uuid
    std.uuid.v4
      @ () -> string
      + returns a new UUID v4
      # identifier
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time

forum
  forum.open
    @ (dsn: string) -> result[forum_state, string]
    + creates a forum bound to a database
    - returns error when the database cannot be opened
    # construction
    -> std.sql.open
  forum.register_user
    @ (state: forum_state, username: string, password: string) -> result[string, string]
    + creates a user and returns its id
    - returns error when the username already exists
    # accounts
    -> std.crypto.bcrypt_hash
    -> std.uuid.v4
    -> std.sql.exec
  forum.authenticate
    @ (state: forum_state, username: string, password: string) -> result[string, string]
    + returns the user id on successful authentication
    - returns error on unknown username
    - returns error on wrong password
    # accounts
    -> std.sql.query
    -> std.crypto.bcrypt_verify
  forum.create_category
    @ (state: forum_state, name: string, description: string) -> result[string, string]
    + creates a category and returns its id
    - returns error when the name already exists
    # content
    -> std.uuid.v4
    -> std.sql.exec
  forum.create_thread
    @ (state: forum_state, user_id: string, category_id: string, title: string, body: string) -> result[string, string]
    + creates a thread with its initial post
    - returns error when the category does not exist
    # content
    -> std.uuid.v4
    -> std.sql.exec
    -> std.time.now_seconds
  forum.reply
    @ (state: forum_state, user_id: string, thread_id: string, body: string) -> result[string, string]
    + creates a post in a thread and returns its id
    - returns error when the thread is locked
    # content
    -> std.uuid.v4
    -> std.sql.exec
    -> std.time.now_seconds
  forum.vote_post
    @ (state: forum_state, user_id: string, post_id: string, value: i32) -> result[void, string]
    + records an up (+1) or down (-1) vote, replacing any prior vote by the user
    - returns error when value is not -1 or +1
    # voting
    -> std.sql.exec
  forum.list_threads
    @ (state: forum_state, category_id: string, limit: i32, offset: i32) -> result[list[thread_summary], string]
    + returns paged threads ordered by last activity descending
    # queries
    -> std.sql.query
  forum.get_thread
    @ (state: forum_state, thread_id: string) -> result[thread_view, string]
    + returns the thread and its posts in order
    - returns error when the thread does not exist
    # queries
    -> std.sql.query
  forum.moderate
    @ (state: forum_state, moderator_id: string, action: string, target_id: string, reason: string) -> result[string, string]
    + records a moderation action and returns its log id
    - returns error on unknown action keyword
    # moderation
    -> std.uuid.v4
    -> std.sql.exec
    -> std.time.now_seconds
  forum.lock_thread
    @ (state: forum_state, moderator_id: string, thread_id: string) -> result[void, string]
    + marks a thread as locked so replies are rejected
    - returns error when the thread does not exist
    # moderation
    -> std.sql.exec
