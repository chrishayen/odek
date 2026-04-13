# Requirement: "an authentication service (register, login, sessions, password reset)"

Five project operations at the feature boundary. The project layer is thin glue; the crypto (bcrypt), token (jwt), and randomness primitives live in std.

std
  std.bcrypt
    std.bcrypt.hash
      @ (password: string) -> result[string, string]
      + returns a bcrypt hash with a random salt
      # cryptography
    std.bcrypt.verify
      @ (password: string, hash: string) -> bool
      + returns true when the password matches the hash
      # cryptography
  std.jwt
    std.jwt.sign
      @ (payload: map[string, string], secret: string) -> result[string, string]
      + signs a payload with HS256
      # token_signing
    std.jwt.verify
      @ (token: string, secret: string) -> result[map[string, string], string]
      + verifies a JWT and returns its payload
      - returns error on bad signature or expired token
      # token_verification
  std.random
    std.random.alphanumeric_string
      @ (length: u32) -> string
      + returns a cryptographically random alphanumeric string of the given length
      # randomness

auth
  auth.register
    @ (username: string, password: string) -> result[user_id, string]
    + hashes the password and creates a new user record
    - returns error when the username is already taken
    # account_management
    -> std.bcrypt.hash
  auth.login
    @ (username: string, password: string) -> result[session_token, string]
    + verifies the password and returns a signed session token
    - returns error on bad password
    - returns error on unknown user
    # account_management
    -> std.bcrypt.verify
    -> std.jwt.sign
  auth.verify_session
    @ (token: session_token) -> result[user_id, string]
    + decodes the session token and returns the authenticated user id
    - returns error on invalid or expired token
    # session
    -> std.jwt.verify
  auth.reset_password_request
    @ (username: string) -> result[reset_token, string]
    + generates a one-time reset token for the user
    - returns error when the username does not exist
    ? reset tokens expire after 1 hour
    # password_reset
    -> std.random.alphanumeric_string
  auth.reset_password_confirm
    @ (reset_token: reset_token, new_password: string) -> result[void, string]
    + verifies the reset token and updates the user's password hash
    - returns error on invalid or expired reset token
    # password_reset
    -> std.bcrypt.hash
