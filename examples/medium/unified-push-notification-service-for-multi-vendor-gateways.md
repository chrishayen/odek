# Requirement: "a unified push notification service that dispatches server-initiated notifications to mobile devices across multiple vendor gateways"

Registers devices under a service key, looks them up at send time, and fans out to the appropriate vendor gateway. Project layer owns registration and routing; std provides http, json, and a kv store abstraction.

std
  std.http
    std.http.post_json
      fn (url: string, headers: map[string, string], body: string) -> result[http_response, string]
      + sends an http POST with a JSON body and returns the response
      - returns error on network failure
      # networking
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization
  std.kv
    std.kv.get
      fn (store: kv_store, key: string) -> result[optional[string], string]
      + returns the value for key, or none when missing
      # storage
    std.kv.put
      fn (store: kv_store, key: string, value: string) -> result[void, string]
      + stores value under key
      # storage
    std.kv.list_prefix
      fn (store: kv_store, prefix: string) -> result[list[string], string]
      + returns all keys starting with prefix
      # storage

push_service
  push_service.new
    fn (store: kv_store) -> service_state
    + wraps a kv store and an empty gateway registry
    # construction
  push_service.register_gateway
    fn (svc: service_state, kind: string, endpoint: string, auth: string) -> service_state
    + registers a vendor gateway under the given kind (e.g. "android", "ios")
    # configuration
  push_service.subscribe_device
    fn (svc: service_state, service_key: string, kind: string, device_token: string) -> result[void, string]
    + stores the device under the service key with its gateway kind
    - returns error when no gateway is registered for kind
    # registration
    -> std.kv.put
  push_service.unsubscribe_device
    fn (svc: service_state, service_key: string, device_token: string) -> result[void, string]
    + removes the device from the service key
    # registration
  push_service.list_devices
    fn (svc: service_state, service_key: string) -> result[list[device_record], string]
    + returns all devices under the service key
    # query
    -> std.kv.list_prefix
    -> std.kv.get
  push_service.send_to_service
    fn (svc: service_state, service_key: string, payload: map[string, string]) -> result[send_report, string]
    + fanouts the payload to every device registered under the service key
    + returns a report with per-device success/failure
    # delivery
    -> std.json.encode_object
    -> std.http.post_json
