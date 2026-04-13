# Requirement: "a message broker inspection and manipulation library"

Library form of a broker swiss-army-knife — exchange/queue introspection plus tap/publish/move operations. No CLI.

std
  std.amqp
    std.amqp.connect
      @ (url: string) -> result[amqp_conn, string]
      + establishes a connection to an AMQP broker
      - returns error on handshake failure
      # transport
    std.amqp.declare_queue
      @ (conn: amqp_conn, name: string, durable: bool) -> result[void, string]
      + declares a queue with the given durability
      - returns error on name collision with incompatible params
      # topology
    std.amqp.publish
      @ (conn: amqp_conn, exchange: string, routing_key: string, body: bytes) -> result[void, string]
      + publishes a single message
      - returns error when the broker nacks the publish
      # publishing
    std.amqp.consume_one
      @ (conn: amqp_conn, queue: string) -> result[optional[amqp_message], string]
      + returns the next message or none if the queue is empty
      # consumption

broker_tools
  broker_tools.list_queues
    @ (mgmt_url: string, credentials: string) -> result[list[queue_info], string]
    + returns queue name, message count, and consumer count
    - returns error when the management endpoint is unreachable
    # introspection
  broker_tools.tap
    @ (conn: amqp_conn, exchange: string, routing_key: string, handler: fn(amqp_message) -> void) -> result[void, string]
    + binds a temporary queue and invokes handler for every matching message until cancelled
    # observation
    -> std.amqp.declare_queue
    -> std.amqp.consume_one
  broker_tools.move_messages
    @ (conn: amqp_conn, from_queue: string, to_exchange: string, to_key: string, limit: i32) -> result[i32, string]
    + moves up to limit messages between destinations and returns the count actually moved
    - returns error when the source queue does not exist
    # routing
    -> std.amqp.consume_one
    -> std.amqp.publish
  broker_tools.publish_bulk
    @ (conn: amqp_conn, exchange: string, routing_key: string, bodies: list[bytes]) -> result[i32, string]
    + publishes every body and returns the count successfully acked
    # publishing
    -> std.amqp.publish
