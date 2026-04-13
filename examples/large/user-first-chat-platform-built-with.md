# Requirement: "a chat platform library with users, rooms, and real-time messaging"

Library-level core: account management, room membership, message persistence, and a subscription fan-out hook. Transport and UI are out of scope.

std
  std.crypto
    std.crypto.random_bytes
      @ (n: i32) -> bytes
      + returns n cryptographically random bytes
      # cryptography
    std.crypto.bcrypt_hash
      @ (password: string, cost: i32) -> result[string, string]
      + returns a bcrypt hash with the given cost factor
      - returns error when cost is outside the valid range
      # cryptography
    std.crypto.bcrypt_verify
      @ (password: string, hash: string) -> bool
      + true when the password matches the hash
      # cryptography
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.id
    std.id.new_uuid
      @ () -> string
      + returns a random v4 UUID string
      # identifiers

chat
  chat.new
    @ () -> chat_state
    + creates an empty chat backend with no users, rooms, or messages
    # construction
  chat.register_user
    @ (state: chat_state, handle: string, password: string) -> result[tuple[user_id, chat_state], string]
    + creates a new user with a hashed password and returns the assigned id
    - returns error when the handle is already taken
    - returns error when the password is empty
    # user_management
    -> std.crypto.bcrypt_hash
    -> std.id.new_uuid
  chat.authenticate
    @ (state: chat_state, handle: string, password: string) -> result[user_id, string]
    + returns the user id when credentials match
    - returns error when the handle is unknown or the password is wrong
    # authentication
    -> std.crypto.bcrypt_verify
  chat.create_room
    @ (state: chat_state, owner: user_id, name: string) -> result[tuple[room_id, chat_state], string]
    + creates a room owned by the user and adds the owner as a member
    - returns error when a room with the same name already exists
    # rooms
    -> std.id.new_uuid
  chat.join_room
    @ (state: chat_state, user: user_id, room: room_id) -> result[chat_state, string]
    + adds the user to the room's member set
    - returns error when the room does not exist
    # membership
  chat.leave_room
    @ (state: chat_state, user: user_id, room: room_id) -> result[chat_state, string]
    + removes the user from the room
    - returns error when the user was not a member
    # membership
  chat.post_message
    @ (state: chat_state, room: room_id, author: user_id, body: string) -> result[tuple[message, chat_state], string]
    + appends a timestamped message to the room's log and returns it
    - returns error when the author is not a member of the room
    - returns error when the body is empty
    # messaging
    -> std.time.now_millis
    -> std.id.new_uuid
  chat.recent_messages
    @ (state: chat_state, room: room_id, limit: i32) -> result[list[message], string]
    + returns up to limit most recent messages ordered oldest-first
    - returns error when the room does not exist
    # history
  chat.subscribe
    @ (state: chat_state, user: user_id, room: room_id, sink: fn(message) -> void) -> result[tuple[subscription_id, chat_state], string]
    + registers a callback invoked for each future message in the room
    - returns error when the user is not a member
    # subscriptions
    -> std.id.new_uuid
  chat.unsubscribe
    @ (state: chat_state, sub: subscription_id) -> chat_state
    + removes the subscription; no-op if the id is unknown
    # subscriptions
