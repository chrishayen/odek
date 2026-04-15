# Requirement: "an in-memory websocket client for unit-testing websocket handlers"

Pairs a handler with a test client so messages can be exchanged without a real network.

std: (all units exist)

wstest
  wstest.new
    fn (handler: ws_handler) -> wstest_client
    + creates a paired client connected to the handler in-memory
    ? no sockets; frames flow through a queue owned by the pair
    # construction
  wstest.send_text
    fn (client: wstest_client, message: string) -> result[wstest_client, string]
    + delivers a text frame to the handler
    - returns error when the connection has been closed
    # client_to_server
  wstest.send_binary
    fn (client: wstest_client, data: bytes) -> result[wstest_client, string]
    + delivers a binary frame to the handler
    - returns error when the connection has been closed
    # client_to_server
  wstest.recv
    fn (client: wstest_client) -> result[ws_frame, string]
    + returns the next frame the handler sent to the client
    - returns error when the connection has been closed and the queue is empty
    # server_to_client
  wstest.close
    fn (client: wstest_client, code: i32, reason: string) -> wstest_client
    + sends a close frame to the handler and marks the client closed
    # teardown
