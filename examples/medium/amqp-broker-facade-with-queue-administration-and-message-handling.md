# Requirement: "an abstraction layer over an AMQP-style message broker for queue administration and publishing and consuming messages"

A thin facade over a broker connection: declare exchanges and queues, bind, publish, and consume. The underlying broker protocol lives in std.

std
  std.amqp
    std.amqp.connect
      fn (url: string) -> result[amqp_conn, string]
      + opens a connection to the broker at the given url
      - returns error on unreachable host or auth failure
      # broker
    std.amqp.close
      fn (conn: amqp_conn) -> void
      + closes the broker connection
      # broker
    std.amqp.declare_exchange
      fn (conn: amqp_conn, name: string, kind: string, durable: bool) -> result[void, string]
      + idempotently declares an exchange with the given kind
      - returns error when a mismatched exchange already exists
      # broker
    std.amqp.declare_queue
      fn (conn: amqp_conn, name: string, durable: bool) -> result[void, string]
      + idempotently declares a queue
      # broker
    std.amqp.bind_queue
      fn (conn: amqp_conn, queue: string, exchange: string, routing_key: string) -> result[void, string]
      + binds the queue to the exchange with the routing key
      # broker
    std.amqp.publish
      fn (conn: amqp_conn, exchange: string, routing_key: string, body: bytes) -> result[void, string]
      + publishes a message via the exchange
      # broker
    std.amqp.consume_one
      fn (conn: amqp_conn, queue: string) -> result[optional[amqp_delivery], string]
      + returns the next delivery or none when the queue is empty
      # broker
    std.amqp.ack
      fn (conn: amqp_conn, delivery: amqp_delivery) -> result[void, string]
      + acknowledges the delivery
      # broker

broker_facade
  broker_facade.open
    fn (url: string) -> result[amqp_conn, string]
    + opens a connection
    - returns error when the broker is unreachable
    # connection
    -> std.amqp.connect
  broker_facade.setup_topic
    fn (conn: amqp_conn, exchange: string, queue: string, routing_key: string) -> result[void, string]
    + declares a durable topic exchange, durable queue, and binding in one call
    - returns error on any declaration failure
    # topology
    -> std.amqp.declare_exchange
    -> std.amqp.declare_queue
    -> std.amqp.bind_queue
  broker_facade.emit
    fn (conn: amqp_conn, exchange: string, routing_key: string, body: bytes) -> result[void, string]
    + publishes a single message
    # publishing
    -> std.amqp.publish
  broker_facade.drain
    fn (conn: amqp_conn, queue: string, max_messages: i32) -> result[list[bytes], string]
    + consumes up to max_messages, acking each, returning their bodies
    + stops early when the queue is empty
    # consuming
    -> std.amqp.consume_one
    -> std.amqp.ack
  broker_facade.shutdown
    fn (conn: amqp_conn) -> void
    + closes the connection
    # connection
    -> std.amqp.close
