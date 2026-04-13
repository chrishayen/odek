# Requirement: "a webhooks gateway for receiving, processing, and reliably delivering messages"

Accepts inbound events, enqueues delivery attempts to downstream endpoints, and retries with exponential backoff.

std
  std.http
    std.http.post
      @ (url: string, headers: map[string, string], body: bytes) -> result[i32, string]
      + performs an HTTP POST and returns the status code
      - returns error on transport failure
      # http
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses any JSON value
      - returns error on malformed input
      # serialization
    std.json.encode_value
      @ (v: json_value) -> string
      + serializes a JSON value
      # serialization
  std.time
    std.time.now_millis
      @ () -> i64
      + returns current unix time in milliseconds
      # time
  std.uuid
    std.uuid.v4
      @ () -> string
      + returns a random UUIDv4 in canonical form
      # identifiers

gateway
  gateway.new
    @ () -> gateway_state
    + creates an empty gateway with no endpoints or queued deliveries
    # construction
  gateway.register_endpoint
    @ (state: gateway_state, event_type: string, url: string) -> gateway_state
    + subscribes a URL to an event type
    # subscriptions
  gateway.receive
    @ (state: gateway_state, event_type: string, payload: json_value) -> tuple[gateway_state, string]
    + enqueues one delivery attempt per registered endpoint and returns the event id
    # ingress
    -> std.uuid.v4
    -> std.time.now_millis
  gateway.next_ready
    @ (state: gateway_state, now_ms: i64) -> optional[delivery_attempt]
    + returns the earliest attempt whose next-try time has passed
    # scheduling
  gateway.attempt
    @ (state: gateway_state, attempt: delivery_attempt) -> tuple[gateway_state, delivery_result]
    + POSTs the payload and records the outcome
    + on 2xx marks the attempt delivered
    - on failure schedules a retry with exponential backoff up to max attempts
    - on exceeded attempts marks the delivery dead
    # delivery
    -> std.http.post
    -> std.json.encode_value
  gateway.retry_delay_ms
    @ (attempt_count: i32) -> i64
    + returns backoff delay in milliseconds for the given attempt number
    + delay doubles per attempt with a cap
    # retry
  gateway.list_dead
    @ (state: gateway_state) -> list[delivery_attempt]
    + returns deliveries that exhausted their retry budget
    # introspection
  gateway.ack_dead
    @ (state: gateway_state, id: string) -> result[gateway_state, string]
    + removes a dead delivery after out-of-band handling
    - returns error when id is not in the dead queue
    # introspection
