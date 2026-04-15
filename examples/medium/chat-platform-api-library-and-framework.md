# Requirement: "a library and framework for a chat platform API"

A typed client for a generic chat platform API. Handles authentication, sends messages, fetches channel history, and dispatches inbound events to registered handlers.

std
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string,string], string]
      + parses a JSON object into a string-keyed map
      - returns error on malformed JSON or non-object root
      # serialization
    std.json.encode_object
      fn (obj: map[string,string]) -> string
      + encodes a string map as a JSON object
      # serialization
  std.http
    std.http.post_json
      fn (url: string, headers: map[string,string], body: string) -> result[string, string]
      + sends a POST and returns the response body on 2xx
      - returns error on non-2xx with the status line included
      # http_client
    std.http.get
      fn (url: string, headers: map[string,string]) -> result[string, string]
      + sends a GET and returns the response body on 2xx
      - returns error on non-2xx
      # http_client

chatapi
  chatapi.new
    fn (token: string, base_url: string) -> client_state
    + creates a client with the given bearer token and API base URL
    - returns a state with an empty token marker when token is empty
    # construction
  chatapi.send_message
    fn (state: client_state, channel_id: string, content: string) -> result[string, string]
    + posts a message and returns the new message id
    - returns error when the channel does not exist or the request fails
    # messaging
    -> std.http.post_json
    -> std.json.encode_object
    -> std.json.parse_object
  chatapi.fetch_history
    fn (state: client_state, channel_id: string, limit: i32) -> result[list[map[string,string]], string]
    + returns up to limit recent messages, newest first
    - returns error when limit is less than 1 or greater than 100
    # history
    -> std.http.get
  chatapi.register_handler
    fn (state: client_state, event_type: string, handler_id: string) -> client_state
    + binds a handler id to an inbound event type
    # event_handlers
  chatapi.dispatch_event
    fn (state: client_state, raw_event: string) -> result[string, string]
    + parses raw_event and returns the handler_id registered for its type
    - returns error when the event has no registered handler
    - returns error when raw_event is not valid JSON
    # event_dispatch
    -> std.json.parse_object
