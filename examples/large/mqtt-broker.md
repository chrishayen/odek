# Requirement: "an embeddable MQTT broker"

A protocol-level MQTT broker: packet codecs, session state, subscription matching, and retained-message storage. The caller owns the network listener.

std: (all units exist)

mqtt
  mqtt.decode_packet
    @ (buffer: bytes) -> result[tuple[mqtt_packet, i32], string]
    + returns a parsed packet and the number of bytes consumed
    - returns error when the buffer is too short for the declared length
    - returns error on unknown packet types
    # codec
  mqtt.encode_packet
    @ (packet: mqtt_packet) -> bytes
    + returns the wire bytes of an MQTT packet
    + encodes remaining-length using variable-byte integers
    # codec
  mqtt.new_broker
    @ () -> broker_state
    + creates an empty broker with no clients and no retained messages
    # construction
  mqtt.connect_client
    @ (state: broker_state, client_id: string, clean_session: bool) -> tuple[broker_state, mqtt_packet]
    + registers the client and returns a CONNACK packet
    + resumes existing subscriptions when clean_session is false
    # lifecycle
  mqtt.disconnect_client
    @ (state: broker_state, client_id: string) -> broker_state
    + removes the client and any non-persistent subscriptions
    # lifecycle
  mqtt.subscribe
    @ (state: broker_state, client_id: string, topic_filter: string, qos: i32) -> tuple[broker_state, mqtt_packet]
    + records the subscription and returns a SUBACK packet
    - the SUBACK reports failure when qos is outside 0..2
    # subscription
  mqtt.unsubscribe
    @ (state: broker_state, client_id: string, topic_filter: string) -> tuple[broker_state, mqtt_packet]
    + removes the matching subscription and returns an UNSUBACK packet
    # subscription
  mqtt.match_topic
    @ (filter: string, topic: string) -> bool
    + returns true when the topic matches an MQTT topic filter
    + "+" matches exactly one level and "#" matches any remaining levels
    - returns false when the filter has a trailing "#" in the middle
    # routing
  mqtt.publish
    @ (state: broker_state, topic: string, payload: bytes, qos: i32, retain: bool) -> tuple[broker_state, list[delivery]]
    + returns the updated state and the list of (client_id, packet) deliveries
    + when retain is true the message replaces any prior retained message for the topic
    # publish
    -> mqtt.match_topic
  mqtt.retained_for
    @ (state: broker_state, topic_filter: string) -> list[tuple[string, bytes]]
    + returns all retained messages matching a topic filter
    # storage
    -> mqtt.match_topic
  mqtt.handle_pingreq
    @ (state: broker_state, client_id: string) -> tuple[broker_state, mqtt_packet]
    + returns a PINGRESP packet and resets the client's keepalive deadline
    # keepalive
