# Requirement: "a client library for a generic chat bot HTTP API"

Wraps a JSON-over-HTTPS bot API: long-polling for updates and typed send methods for replies.

std
  std.http
    std.http.get
      fn (url: string) -> result[http_response, string]
      + performs an HTTPS GET and returns status and body
      - returns error on network failure
      # http
    std.http.post_json
      fn (url: string, body: string) -> result[http_response, string]
      + performs an HTTPS POST with the given JSON body
      - returns error on network failure
      # http
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses a JSON document into a dynamic value
      - returns error on malformed JSON
      # serialization
    std.json.encode
      fn (value: json_value) -> string
      + encodes a dynamic value as a JSON string
      # serialization
    std.json.get_field
      fn (value: json_value, path: string) -> optional[json_value]
      + returns the nested field referenced by a dot-separated path or none
      # serialization
  std.url
    std.url.encode_query
      fn (params: map[string, string]) -> string
      + builds a percent-encoded query string from the parameter map
      # url

chat_bot
  chat_bot.new
    fn (api_base: string, token: string) -> bot_state
    + creates a bot client targeting the given API base and authorization token
    # construction
  chat_bot.get_me
    fn (bot: bot_state) -> result[bot_identity, string]
    + returns the bot's own identity from the remote API
    - returns error when the token is invalid
    # identity
    -> std.http.get
    -> std.json.parse
  chat_bot.poll_updates
    fn (bot: bot_state, offset: i64, timeout_sec: i32) -> result[tuple[list[update], i64], string]
    + long-polls for updates and returns them with the next offset to use
    - returns error when the response is not a success envelope
    # polling
    -> std.http.get
    -> std.url.encode_query
    -> std.json.parse
  chat_bot.send_message
    fn (bot: bot_state, chat_id: i64, text: string) -> result[message_id, string]
    + sends a plain text message to the given chat
    - returns error when the chat does not exist
    # sending
    -> std.http.post_json
    -> std.json.encode
    -> std.json.parse
  chat_bot.send_reply
    fn (bot: bot_state, chat_id: i64, reply_to: message_id, text: string) -> result[message_id, string]
    + sends a message as a reply to an existing message
    # sending
    -> std.http.post_json
    -> std.json.encode
  chat_bot.edit_message
    fn (bot: bot_state, chat_id: i64, mid: message_id, text: string) -> result[void, string]
    + replaces the text of a previously sent message
    - returns error when the message is too old to edit
    # sending
    -> std.http.post_json
    -> std.json.encode
  chat_bot.delete_message
    fn (bot: bot_state, chat_id: i64, mid: message_id) -> result[void, string]
    + deletes a previously sent message
    # sending
    -> std.http.post_json
  chat_bot.answer_callback
    fn (bot: bot_state, callback_id: string, text: string) -> result[void, string]
    + acknowledges an inline-button callback with an optional notification
    # sending
    -> std.http.post_json
    -> std.json.encode
  chat_bot.parse_update
    fn (raw: string) -> result[update, string]
    + decodes a single update envelope from its JSON representation
    - returns error when required fields are missing
    # parsing
    -> std.json.parse
    -> std.json.get_field
