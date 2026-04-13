# Requirement: "a chat application backend"

The flagship large example. Six project operations at the feature boundary; all the substantive plumbing — bcrypt, JWT, WebSocket, SQL — lives in std as genuinely reusable subsystems.

std
  std.bcrypt
    std.bcrypt.hash
      @ (password: string) -> result[string, string]
      + returns a bcrypt hash of the password with a random salt
      - returns error when password exceeds 72 bytes
      # cryptography
    std.bcrypt.verify
      @ (password: string, hash: string) -> bool
      + returns true when the password matches the hash
      + returns false on mismatch
      # cryptography
  std.jwt
    std.jwt.sign
      @ (payload: map[string, string], secret: string) -> result[string, string]
      + signs a payload with HS256 and returns a JWT
      - returns error when secret is empty
      # token_signing
    std.jwt.verify
      @ (token: string, secret: string) -> result[map[string, string], string]
      + verifies a JWT and returns its payload
      - returns error on bad signature or expired token
      # token_verification
  std.websocket
    std.websocket.upgrade
      @ (req: http_request) -> result[websocket_conn, string]
      + upgrades an HTTP request to a WebSocket connection
      - returns error when upgrade headers are missing
      # networking
    std.websocket.send
      @ (c: websocket_conn, data: bytes) -> result[void, string]
      + sends a binary or text frame on the connection
      - returns error when the connection is closed
      # networking
  std.sql
    std.sql.query
      @ (conn: db_conn, sql: string, args: list[any]) -> result[list[row], string]
      + executes a parameterized query and returns the result rows
      - returns error on syntax error
      - returns error on constraint violation
      # persistence

chat
  chat.create_user
    @ (creds: credentials) -> result[user_id, string]
    + registers a new user with the password stored as a bcrypt hash
    - returns error when the username is already taken
    # account_management
    -> std.bcrypt.hash
    -> std.sql.query
  chat.authenticate
    @ (creds: credentials) -> result[session_token, string]
    + verifies the password and returns a signed session token
    - returns error on bad password
    - returns error on unknown user
    # account_management
    -> std.bcrypt.verify
    -> std.jwt.sign
    -> std.sql.query
  chat.create_room
    @ (token: session_token, name: string) -> result[room_id, string]
    + creates a new chat room owned by the authenticated user
    - returns error when the token is invalid
    # room_management
    -> std.jwt.verify
    -> std.sql.query
  chat.join_room
    @ (token: session_token, room_id: room_id) -> result[void, string]
    + adds the authenticated caller to the room
    - returns error when the token is invalid
    # room_management
    -> std.jwt.verify
    -> std.sql.query
  chat.send_message
    @ (token: session_token, room_id: room_id, body: string) -> result[message_id, string]
    + stores the message and broadcasts it to connected room members
    - returns error when the user is not a member of the room
    # messaging
    -> std.jwt.verify
    -> std.sql.query
    -> std.websocket.send
  chat.fetch_messages
    @ (token: session_token, room_id: room_id, since: i64) -> result[list[message], string]
    + returns messages posted to the room since the given unix timestamp
    - returns error when the user is not a member of the room
    # messaging
    -> std.jwt.verify
    -> std.sql.query
