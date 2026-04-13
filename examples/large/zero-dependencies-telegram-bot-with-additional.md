# Requirement: "a chat bot library with message handling and interactive UI components"

Handles long-polling updates from a chat platform API, dispatches to registered handlers, and builds interactive UI elements.

std
  std.http
    std.http.get
      @ (url: string) -> result[http_response, string]
      + performs an HTTP GET request
      - returns error on transport failure or non-2xx status
      # networking
    std.http.post_json
      @ (url: string, body: string) -> result[http_response, string]
      + performs an HTTP POST with a JSON body
      - returns error on transport failure
      # networking
  std.json
    std.json.parse
      @ (raw: string) -> result[json_value, string]
      + parses a JSON document into a generic value
      - returns error on malformed input
      # serialization
    std.json.encode
      @ (value: json_value) -> string
      + encodes a generic value to a JSON string
      # serialization

chatbot
  chatbot.new_client
    @ (api_base: string, token: string) -> client_state
    + creates a client bound to the given API base URL and auth token
    # construction
  chatbot.fetch_updates
    @ (state: client_state, since_id: i64) -> result[tuple[list[update], i64], string]
    + fetches updates newer than since_id and returns them with the new high-water mark
    - returns error on transport or parse failure
    # polling
    -> std.http.get
    -> std.json.parse
  chatbot.send_message
    @ (state: client_state, chat_id: i64, text: string) -> result[i64, string]
    + sends a text message and returns the resulting message id
    - returns error when the API rejects the request
    # messaging
    -> std.json.encode
    -> std.http.post_json
  chatbot.send_with_keyboard
    @ (state: client_state, chat_id: i64, text: string, kb: keyboard) -> result[i64, string]
    + sends a message accompanied by an inline keyboard
    - returns error on transport failure
    # messaging
    -> std.json.encode
    -> std.http.post_json
  chatbot.new_router
    @ () -> router_state
    + creates an empty update router
    # routing
  chatbot.on_command
    @ (router: router_state, command: string, handler: handler_id) -> router_state
    + registers a handler for a slash-prefixed command
    # routing
  chatbot.on_callback
    @ (router: router_state, data_prefix: string, handler: handler_id) -> router_state
    + registers a handler for callback-button events matching a prefix
    # routing
  chatbot.dispatch
    @ (router: router_state, upd: update) -> list[handler_id]
    + returns the handler ids that should process the update
    # routing
  chatbot.build_keyboard
    @ (rows: list[list[button]]) -> keyboard
    + constructs a keyboard from a matrix of buttons
    # ui
  chatbot.button_text
    @ (label: string, callback_data: string) -> button
    + creates a button with a label and a callback payload
    # ui
  chatbot.button_url
    @ (label: string, url: string) -> button
    + creates a button that opens a URL when pressed
    # ui
  chatbot.edit_message
    @ (state: client_state, chat_id: i64, message_id: i64, text: string) -> result[void, string]
    + replaces the text of a previously sent message
    - returns error when the API rejects the edit
    # messaging
    -> std.json.encode
    -> std.http.post_json
