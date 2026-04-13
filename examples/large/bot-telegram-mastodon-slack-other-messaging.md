# Requirement: "a messaging bot that archives web pages on command"

A chat bot with pluggable platform adapters, a command router, and an archive pipeline that snapshots URLs.

std
  std.http
    std.http.get
      @ (url: string) -> result[http_response, string]
      + performs an HTTP GET and returns status, headers, and body
      - returns error on DNS or connection failure
      # networking
    std.http.post
      @ (url: string, headers: map[string, string], body: bytes) -> result[http_response, string]
      + performs an HTTP POST and returns status, headers, and body
      - returns error on DNS or connection failure
      # networking
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns the current unix time in seconds
      # time
  std.encoding
    std.encoding.url_encode
      @ (raw: string) -> string
      + percent-encodes reserved characters
      # encoding
  std.hash
    std.hash.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex SHA-256 digest
      # hashing

wayback
  wayback.new_bot
    @ () -> bot_state
    + creates an empty bot with no adapters and no commands
    # construction
  wayback.register_adapter
    @ (bot: bot_state, name: string, adapter: platform_adapter) -> bot_state
    + attaches a chat platform adapter identified by name
    # platform
  wayback.register_command
    @ (bot: bot_state, trigger: string, handler: fn(command_ctx) -> command_response) -> bot_state
    + binds a command trigger (e.g. "/archive") to a handler
    # routing
  wayback.route_message
    @ (bot: bot_state, incoming: incoming_message) -> optional[command_response]
    + parses the message, dispatches to the matching command handler, and returns the response
    - returns none when no command trigger is present
    # routing
  wayback.extract_urls
    @ (text: string) -> list[string]
    + returns every http(s) URL found in the message text in order
    - returns an empty list when no URL is present
    # parsing
  wayback.archive_url
    @ (url: string, sink: archive_sink) -> result[archive_record, string]
    + fetches the URL, computes a content hash, and hands the body to the archive sink
    - returns error when the fetch fails
    - returns error when the response status is not 2xx
    # archiving
    -> std.http.get
    -> std.hash.sha256_hex
    -> std.time.now_seconds
  wayback.archive_handler
    @ (ctx: command_ctx, sink: archive_sink) -> command_response
    + handles the /archive command by archiving every URL in the message and replying with results
    # commands
    -> wayback.extract_urls
    -> wayback.archive_url
  wayback.format_reply
    @ (records: list[archive_record], failures: list[string]) -> string
    + renders a human-readable reply summarizing archived URLs and failures
    # formatting
  wayback.send_reply
    @ (bot: bot_state, adapter_name: string, channel: string, text: string) -> result[void, string]
    + dispatches a reply through the named adapter
    - returns error when the adapter is not registered
    # platform
