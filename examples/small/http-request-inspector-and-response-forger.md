# Requirement: "a tool for inspecting incoming http requests and forging responses"

Captures requests in a buffer and lets callers define the response the next request will receive.

std: (all units exist)

http_inspect
  http_inspect.new
    @ () -> inspect_state
    + returns an inspector with no captured requests and a default 200 empty response
    # construction
  http_inspect.set_response
    @ (state: inspect_state, status: i32, headers: map[string,string], body: bytes) -> inspect_state
    + installs the response the next incoming request will receive
    # configuration
  http_inspect.handle
    @ (state: inspect_state, req: http_request) -> tuple[http_response, inspect_state]
    + records the request in the buffer and returns the configured response
    # dispatch
  http_inspect.captures
    @ (state: inspect_state) -> list[http_request]
    + returns all recorded requests in arrival order
    # inspection
