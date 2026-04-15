# Requirement: "a chat bot framework"

Bot framework for a generic chat platform: registers commands, dispatches incoming messages, and sends replies through a pluggable transport.

std
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.http
    std.http.post_json
      fn (url: string, headers: map[string, string], body: string) -> result[string, string]
      + performs an HTTP POST with a JSON body and returns the response body
      - returns error on non-2xx status or network failure
      # http_client
  std.strings
    std.strings.split
      fn (input: string, sep: string) -> list[string]
      + splits input on the separator preserving order
      + returns a single-element list when sep is absent
      # strings
    std.strings.trim
      fn (input: string) -> string
      + removes leading and trailing whitespace
      # strings

chatbot
  chatbot.new
    fn (prefix: string) -> bot_state
    + creates a bot that recognizes commands starting with the prefix
    # construction
  chatbot.register_command
    fn (state: bot_state, name: string, handler: fn(args: list[string], ctx: message_context) -> string) -> result[bot_state, string]
    + adds a handler keyed by command name
    - returns error when name is already registered
    # registration
  chatbot.parse_message
    fn (state: bot_state, raw: string) -> optional[parsed_command]
    + returns command name and args when the text begins with the prefix
    - returns empty when text does not start with the prefix
    # parsing
    -> std.strings.split
    -> std.strings.trim
  chatbot.dispatch
    fn (state: bot_state, cmd: parsed_command, ctx: message_context) -> result[string, string]
    + invokes the registered handler and returns its reply text
    - returns error when command name is unknown
    # dispatch
  chatbot.handle_event
    fn (state: bot_state, raw_event: string) -> result[optional[string], string]
    + parses an incoming event and dispatches to a handler when it is a command
    - returns error on malformed event payload
    # event_loop
    -> std.json.parse_object
  chatbot.send_reply
    fn (state: bot_state, channel: string, text: string) -> result[void, string]
    + posts a reply to the given channel via the configured transport
    - returns error when the transport rejects the request
    # transport
    -> std.json.encode_object
    -> std.http.post_json
  chatbot.add_middleware
    fn (state: bot_state, middleware: fn(ctx: message_context, next: fn() -> string) -> string) -> bot_state
    + appends a middleware that wraps every dispatched handler
    # middleware
  chatbot.set_transport
    fn (state: bot_state, transport: transport_config) -> bot_state
    + replaces the outbound message transport
    ? transport is pluggable so tests can substitute a fake sink
    # configuration
