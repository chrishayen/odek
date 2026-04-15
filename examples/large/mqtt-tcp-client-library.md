# Requirement: "a client library for a pub-sub messaging protocol over TCP"

Connect, subscribe to topics, publish messages, and handle the wire-level framing.

std
  std.net
    std.net.tcp_connect
      fn (host: string, port: u16) -> result[tcp_conn, string]
      + returns a connected socket
      - returns error on dns or connection failure
      # network
    std.net.tcp_write
      fn (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes all bytes
      - returns error on broken pipe
      # network
    std.net.tcp_read
      fn (conn: tcp_conn, max: i32) -> result[bytes, string]
      + reads up to max bytes
      - returns error on eof
      # network
    std.net.tcp_close
      fn (conn: tcp_conn) -> void
      + closes the socket
      # network
  std.encoding
    std.encoding.varint_encode
      fn (value: u32) -> bytes
      + encodes an unsigned integer as a length-prefixed varint
      # encoding
    std.encoding.varint_decode
      fn (data: bytes, offset: i32) -> result[tuple[u32, i32], string]
      + returns (value, bytes_consumed)
      - returns error on truncated input
      # encoding

mqtt
  mqtt.encode_connect
    fn (client_id: string, keepalive: u16, clean_session: bool) -> bytes
    + returns the framed CONNECT packet
    # packet
    -> std.encoding.varint_encode
  mqtt.encode_publish
    fn (topic: string, payload: bytes, qos: u8, packet_id: u16) -> bytes
    + returns the framed PUBLISH packet
    # packet
    -> std.encoding.varint_encode
  mqtt.encode_subscribe
    fn (packet_id: u16, topic_filters: list[tuple[string, u8]]) -> bytes
    + returns the framed SUBSCRIBE packet
    # packet
    -> std.encoding.varint_encode
  mqtt.decode_packet
    fn (data: bytes) -> result[tuple[mqtt_packet, i32], string]
    + returns the parsed packet and the number of bytes consumed
    - returns error on truncated or unknown packet type
    # packet
    -> std.encoding.varint_decode
  mqtt.connect
    fn (host: string, port: u16, client_id: string) -> result[mqtt_client, string]
    + dials the server, sends CONNECT, and waits for CONNACK
    - returns error on rejected credentials
    # session
    -> std.net.tcp_connect
    -> std.net.tcp_write
    -> std.net.tcp_read
  mqtt.publish
    fn (client: mqtt_client, topic: string, payload: bytes, qos: u8) -> result[void, string]
    + sends a PUBLISH and waits for ack when qos > 0
    - returns error when the connection is closed
    # publish
    -> std.net.tcp_write
  mqtt.subscribe
    fn (client: mqtt_client, topic_filters: list[tuple[string, u8]]) -> result[list[u8], string]
    + sends SUBSCRIBE and returns granted qos per filter
    - returns error on SUBACK failure codes
    # subscribe
    -> std.net.tcp_write
    -> std.net.tcp_read
  mqtt.next_message
    fn (client: mqtt_client) -> result[mqtt_message, string]
    + returns the next published message delivered to a subscribed topic
    - returns error when the connection is closed
    # receive
    -> std.net.tcp_read
  mqtt.disconnect
    fn (client: mqtt_client) -> void
    + sends DISCONNECT and closes the socket
    # session
    -> std.net.tcp_write
    -> std.net.tcp_close
  mqtt.topic_matches
    fn (filter: string, topic: string) -> bool
    + returns true when the filter matches the topic including "+" and "#" wildcards
    # routing
