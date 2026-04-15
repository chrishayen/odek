# Requirement: "a client library for a messaging platform bot API with type-safe models, automatic retries, and rate-limit handling"

Typed wrappers over a JSON/HTTP bot API. Retries and rate-limiting sit in front of every call.

std
  std.http
    std.http.request
      fn (method: string, url: string, headers: map[string,string], body: bytes) -> result[http_response, string]
      + performs an HTTP request and returns status, headers, and body
      - returns error when the transport fails
      # networking
  std.json
    std.json.encode
      fn (value: json_value) -> string
      + encodes a JSON value as a string
      # serialization
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON value
      - returns error on malformed input
      # serialization
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
    std.time.sleep_millis
      fn (ms: i64) -> void
      + suspends the current task for ms milliseconds
      # time

bot
  bot.new
    fn (token: string, base_url: string) -> bot_client
    + creates a client with the given credentials and endpoint
    - panics or errors when token is empty
    # construction
  bot.with_retry
    fn (client: bot_client, max_attempts: i32, backoff_ms: i64) -> bot_client
    + configures exponential retry on transient failures
    # resilience
  bot.with_rate_limit
    fn (client: bot_client, per_second: f64) -> bot_client
    + adds a token-bucket limiter applied before every call
    # resilience
    -> std.time.now_millis
  bot.call
    fn (client: bot_client, method: string, payload: json_value) -> result[json_value, string]
    + performs an API call, waiting on the limiter and retrying transient errors
    - returns error on non-retryable API failure
    # request
    -> std.http.request
    -> std.json.encode
    -> std.json.parse
    -> std.time.sleep_millis
  bot.send_message
    fn (client: bot_client, chat_id: string, text: string) -> result[message, string]
    + returns the posted message record
    - returns error when the chat does not exist
    # messaging
    -> bot.call
  bot.edit_message
    fn (client: bot_client, chat_id: string, message_id: string, text: string) -> result[message, string]
    + returns the updated message record
    - returns error when the message is not editable
    # messaging
    -> bot.call
  bot.delete_message
    fn (client: bot_client, chat_id: string, message_id: string) -> result[bool, string]
    + returns true when the message was removed
    - returns error when the message does not exist
    # messaging
    -> bot.call
  bot.get_updates
    fn (client: bot_client, offset: i64, timeout_secs: i32) -> result[list[update], string]
    + returns updates newer than offset using long polling
    # polling
    -> bot.call
  bot.answer_callback
    fn (client: bot_client, callback_id: string, text: string) -> result[bool, string]
    + acknowledges an inline callback
    # messaging
    -> bot.call
