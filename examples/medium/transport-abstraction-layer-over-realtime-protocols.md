# Requirement: "an abstraction layer over real-time transport frameworks so callers are not locked into one implementation"

The library defines a transport interface and a connection manager on top. Concrete transports (websocket, long-poll, etc.) are injected.

std: (all units exist)

transport
  transport.register_driver
    @ (name: string, open_fn: driver_open_fn, send_fn: driver_send_fn, close_fn: driver_close_fn) -> result[void, string]
    + registers a named transport driver in the global registry
    - returns error when a driver with the same name already exists
    # registry
  transport.open
    @ (driver_name: string, endpoint: string) -> result[connection, string]
    + opens a connection using the named driver
    - returns error when the driver is unknown
    - returns error when the driver's open function fails
    # lifecycle
    -> transport.register_driver
  transport.send
    @ (conn: connection, payload: bytes) -> result[void, string]
    + sends a payload through the underlying driver
    - returns error when the connection is already closed
    # messaging
  transport.on_message
    @ (conn: connection, handler: fn(bytes) -> void) -> connection
    + registers a callback invoked for each incoming payload
    ? replaces any previously registered handler
    # messaging
  transport.close
    @ (conn: connection) -> result[void, string]
    + closes the connection and releases driver resources
    + safe to call more than once; second call is a no-op
    # lifecycle
