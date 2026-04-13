# Requirement: "a socket mock framework"

Registers canned responses for host:port pairs so tests can run without a real network.

std: (all units exist)

socket_mock
  socket_mock.new
    @ () -> mock_registry
    + creates an empty registry with no entries
    # construction
  socket_mock.register
    @ (reg: mock_registry, host: string, port: i32, script: list[mock_exchange]) -> void
    + binds a scripted sequence of (expected_send, canned_reply) exchanges to a host:port
    # registration
  socket_mock.dial
    @ (reg: mock_registry, host: string, port: i32) -> result[mock_conn, string]
    + returns a mock connection bound to the registered script
    - returns error when no script is registered for the target
    # connection
  socket_mock.send
    @ (c: mock_conn, data: bytes) -> result[void, string]
    + consumes the next expected_send from the script
    - returns error when the script has been exhausted
    - returns error when the data does not match the next expected send
    # interaction
  socket_mock.recv
    @ (c: mock_conn) -> result[bytes, string]
    + returns the canned reply that pairs with the most recent send
    - returns error when no reply is pending
    # interaction
