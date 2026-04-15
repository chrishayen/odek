# Requirement: "a CoAP client and server implementation (RFC 7252)"

A library that encodes and decodes CoAP messages and tracks request/response state for both client and server roles. UDP I/O is the host's responsibility.

std
  std.time
    std.time.now_millis
      fn () -> i64
      + returns current unix time in milliseconds
      # time
  std.random
    std.random.u16
      fn () -> u16
      + returns a random 16-bit integer
      # random

coap
  coap.encode_message
    fn (msg: coap_message) -> bytes
    + serializes a CoAP message to its binary wire format per RFC 7252
    ? uses delta-encoded options as specified by the RFC
    # encoding
  coap.decode_message
    fn (data: bytes) -> result[coap_message, string]
    + parses bytes into a CoAP message
    - returns error when the version field is not 1
    - returns error when option deltas exceed the buffer length
    # decoding
  coap.new_client
    fn () -> client_state
    + creates a client tracking outstanding requests by token
    # construction
  coap.build_request
    fn (state: client_state, method: string, uri_path: string, payload: bytes) -> tuple[bytes, client_state]
    + builds a confirmable request and records it so the matching response can be dispatched
    - the method must be one of GET, POST, PUT, DELETE
    # request
    -> std.random.u16
  coap.on_client_datagram
    fn (state: client_state, data: bytes) -> result[tuple[coap_message, client_state], string]
    + decodes an incoming datagram and matches it against a pending request by token
    - returns error when no pending request matches the token
    # response
  coap.new_server
    fn () -> server_state
    + creates a server with no registered resources
    # construction
  coap.register_resource
    fn (state: server_state, uri_path: string, handler_name: string) -> server_state
    + registers a named handler for a URI path
    ? handler_name is resolved externally; this only records the mapping
    # routing
  coap.dispatch_request
    fn (state: server_state, data: bytes) -> result[tuple[string, coap_message], string]
    + decodes a request datagram and returns (handler_name, request_message)
    - returns error when no resource matches the URI
    - returns error when the message cannot be decoded
    # routing
  coap.build_response
    fn (request: coap_message, code: string, payload: bytes) -> bytes
    + builds an acknowledgment response matching the request's token and id
    # response
  coap.retransmit_due
    fn (state: client_state) -> tuple[list[bytes], client_state]
    + returns datagrams whose retransmission timer has elapsed
    ? uses the exponential backoff specified in RFC 7252 section 4.2
    # reliability
    -> std.time.now_millis
