# Requirement: "a web framework whose transports and codecs are pluggable"

A handler is registered once and can be served over multiple transports (HTTP, WebSocket, etc.) with multiple encodings (JSON, binary, etc.). The framework dispatches based on transport-provided messages.

std: (all units exist)

pluggable_framework
  pluggable_framework.new
    fn () -> framework_state
    + returns an empty framework
    # construction
  pluggable_framework.register_transport
    fn (state: framework_state, name: string, transport: transport_impl) -> framework_state
    + adds a named transport
    # registration
  pluggable_framework.register_codec
    fn (state: framework_state, name: string, codec: codec_impl) -> framework_state
    + adds a named codec
    # registration
  pluggable_framework.register_handler
    fn (state: framework_state, route: string, handler: handler_fn) -> framework_state
    + binds a handler to a route
    # registration
  pluggable_framework.handle_message
    fn (state: framework_state, transport: string, codec: string, raw: bytes) -> result[bytes, string]
    + decodes the message, dispatches to the handler, and encodes the reply
    - returns error when the transport or codec is unknown
    - returns error when no handler matches the route
    # dispatch
  pluggable_framework.list_routes
    fn (state: framework_state) -> list[string]
    + returns all registered routes
    # inspection
