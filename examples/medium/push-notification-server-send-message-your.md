# Requirement: "a push notification server that sends messages to mobile devices"

Core library: register device tokens, enqueue messages, dispatch to a pluggable delivery backend. Transport details live behind a std primitive.

std
  std.time
    std.time.now_seconds
      @ () -> i64
      + returns current unix time in seconds
      # time
  std.uuid
    std.uuid.generate
      @ () -> string
      + returns a random uuid string
      # identifiers
  std.json
    std.json.encode_object
      @ (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

push_server
  push_server.new
    @ () -> push_state
    + creates an empty push server state with no registered devices
    # construction
  push_server.register_device
    @ (state: push_state, device_token: string, user_id: string) -> result[push_state, string]
    + adds a device token associated with a user and returns updated state
    - returns error when device_token is empty
    # device_registry
  push_server.unregister_device
    @ (state: push_state, device_token: string) -> push_state
    + removes the device token if present and returns updated state
    + returns unchanged state when token is unknown
    # device_registry
  push_server.enqueue_message
    @ (state: push_state, user_id: string, title: string, body: string) -> tuple[push_state, string]
    + appends a message for every device belonging to user_id and returns (state, message_id)
    ? message_id is generated so callers can track delivery
    # message_queue
    -> std.uuid.generate
    -> std.time.now_seconds
  push_server.pending_for_device
    @ (state: push_state, device_token: string) -> list[string]
    + returns JSON-encoded pending messages for a specific device in enqueue order
    # message_queue
    -> std.json.encode_object
  push_server.mark_delivered
    @ (state: push_state, message_id: string) -> push_state
    + removes the message from all device queues and returns updated state
    # delivery_tracking
