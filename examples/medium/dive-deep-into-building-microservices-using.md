# Requirement: "a microservices toolkit for wiring services with typed rpc calls"

A small surface for defining services, routing incoming calls to handlers, and invoking remote procedures over a pluggable transport.

std
  std.encoding
    std.encoding.encode_request
      @ (method: string, body: bytes) -> bytes
      + encodes a method name and body into a framed request
      # serialization
    std.encoding.decode_request
      @ (frame: bytes) -> result[tuple[string, bytes], string]
      + decodes a framed request into (method, body)
      - returns error on truncated input
      # serialization
    std.encoding.encode_response
      @ (ok: bool, body: bytes) -> bytes
      + encodes a response with a status flag and body
      # serialization
    std.encoding.decode_response
      @ (frame: bytes) -> result[tuple[bool, bytes], string]
      + decodes a response frame into (ok, body)
      - returns error on truncated input
      # serialization

microservices
  microservices.new_service
    @ (name: string) -> service_state
    + creates a service with the given name and an empty method table
    # construction
  microservices.register_method
    @ (state: service_state, method: string) -> service_state
    + records that the service exposes a named method
    + replaces an existing registration
    # service_definition
  microservices.handle_request
    @ (state: service_state, frame: bytes) -> bytes
    + decodes a request, dispatches to the registered method, and returns an encoded response
    + returns an error response when the method is not registered
    # dispatch
    -> std.encoding.decode_request
    -> std.encoding.encode_response
  microservices.new_client
    @ (target: string) -> client_state
    + creates a client bound to a transport target
    # construction
  microservices.call
    @ (client: client_state, method: string, body: bytes) -> result[bytes, string]
    + sends an encoded request, waits for the response, and returns the body
    - returns error when the response status flag is false
    # rpc
    -> std.encoding.encode_request
    -> std.encoding.decode_response
