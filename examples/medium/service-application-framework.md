# Requirement: "a framework for building applications and services"

A minimal service framework: register named handlers on routes, dispatch incoming requests, and return structured responses. The transport layer is out of scope; the framework focuses on routing and handler lifecycle.

std
  std.collections
    std.collections.map_get_or
      @ (m: map[string,string], key: string, fallback: string) -> string
      + returns the value when key is present
      + returns fallback when key is absent
      # collections

service
  service.new
    @ () -> service_state
    + creates an empty service with no registered routes
    # construction
  service.register
    @ (state: service_state, method: string, path: string, handler_id: string) -> service_state
    + stores a (method, path) -> handler_id mapping
    + later registrations for the same key overwrite earlier ones
    # registration
  service.dispatch
    @ (state: service_state, method: string, path: string) -> result[string, string]
    + returns the handler_id registered for (method, path)
    - returns error "not found" when no matching route exists
    # routing
  service.handle_request
    @ (state: service_state, method: string, path: string, headers: map[string,string], body: bytes) -> response
    + wraps dispatch and builds a response record with status, headers, and body
    - response status is 404 when dispatch returns an error
    # request_lifecycle
    -> std.collections.map_get_or
  service.shutdown
    @ (state: service_state) -> void
    + releases any resources held by registered handlers
    ? handlers are owned by the caller; shutdown only clears internal tables
    # lifecycle
