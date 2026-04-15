# Requirement: "a chat bot framework with pluggable commands and external service integrations"

A chat bot connects to a chat service, dispatches incoming messages to registered commands, and exposes hooks for external service integrations.

std
  std.http
    std.http.send
      fn (method: string, url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs the request and returns status, headers, and body
      - returns error on network failure
      # http
  std.ws
    std.ws.connect
      fn (url: string) -> result[ws_conn, string]
      + opens a websocket connection to url
      - returns error when the handshake fails
      # websocket
    std.ws.receive
      fn (conn: ws_conn) -> result[string, string]
      + returns the next text message
      - returns error when the connection is closed
      # websocket
    std.ws.send
      fn (conn: ws_conn, message: string) -> result[void, string]
      + sends a text message
      - returns error when the connection is closed
      # websocket
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses JSON text into a value tree
      - returns error on invalid JSON
      # serialization
    std.json.encode
      fn (value: json_value) -> string
      + serializes a JSON value to a string
      # serialization

chatbot
  chatbot.new
    fn (token: string, service_url: string) -> chatbot_state
    + creates a bot with credentials for the chat service
    # construction
  chatbot.register_command
    fn (bot: chatbot_state, name: string, handler: command_handler) -> chatbot_state
    + returns a bot with the named command registered
    ? commands are matched when the incoming text begins with the name
    # registration
  chatbot.register_integration
    fn (bot: chatbot_state, name: string, fetcher: integration_fetcher) -> chatbot_state
    + returns a bot with an external integration registered by name
    # registration
  chatbot.connect
    fn (bot: chatbot_state) -> result[chatbot_session, string]
    + opens a real-time connection to the chat service
    - returns error when authentication fails
    # connection
    -> std.http.send
    -> std.ws.connect
  chatbot.parse_event
    fn (raw: string) -> result[chat_event, string]
    + decodes a chat event from its wire format
    - returns error when the payload is not a known event type
    # parsing
    -> std.json.parse
  chatbot.dispatch
    fn (bot: chatbot_state, event: chat_event) -> optional[command_result]
    + routes an event to the matching command and returns its result
    - returns none when no command matches
    # dispatch
  chatbot.call_integration
    fn (bot: chatbot_state, name: string, query: map[string, string]) -> result[json_value, string]
    + invokes a named integration's fetcher
    - returns error when no integration with the name is registered
    # integration
  chatbot.send_reply
    fn (session: chatbot_session, channel: string, text: string) -> result[void, string]
    + posts a reply in the given channel
    - returns error when the connection is closed
    # messaging
    -> std.json.encode
    -> std.ws.send
  chatbot.run
    fn (bot: chatbot_state, session: chatbot_session) -> result[void, string]
    + reads events, dispatches to commands, and sends replies until the session ends
    - returns error when the connection drops unexpectedly
    # event_loop
    -> std.ws.receive
