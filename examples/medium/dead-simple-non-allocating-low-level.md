# Requirement: "a low-level MQTT protocol codec suitable for constrained environments"

Encode and decode MQTT control packets against caller-provided buffers. No networking.

std: (all units exist)

mqtt
  mqtt.encode_varint
    @ (value: u32, out: bytes) -> result[i32, string]
    + writes the MQTT remaining-length varint to out and returns bytes written
    - returns error when out is shorter than required
    # wire
  mqtt.decode_varint
    @ (src: bytes) -> result[tuple[u32, i32], string]
    + returns the decoded value and number of bytes consumed
    - returns error when the value exceeds four bytes
    # wire
  mqtt.encode_connect
    @ (client_id: string, keepalive: u16, out: bytes) -> result[i32, string]
    + writes a CONNECT packet and returns bytes written
    - returns error when out is too small
    # packet
    -> mqtt.encode_varint
  mqtt.encode_publish
    @ (topic: string, payload: bytes, qos: u8, out: bytes) -> result[i32, string]
    + writes a PUBLISH packet and returns bytes written
    - returns error when qos exceeds 2
    # packet
    -> mqtt.encode_varint
  mqtt.encode_subscribe
    @ (packet_id: u16, topic: string, qos: u8, out: bytes) -> result[i32, string]
    + writes a SUBSCRIBE packet with one topic filter
    - returns error when out is too small
    # packet
    -> mqtt.encode_varint
  mqtt.decode_fixed_header
    @ (src: bytes) -> result[tuple[u8, u32, i32], string]
    + returns packet type, remaining length, and header size
    - returns error when src is shorter than the header
    # packet
    -> mqtt.decode_varint
  mqtt.decode_publish
    @ (src: bytes) -> result[publish_packet, string]
    + returns topic, qos, packet id, and payload slice
    - returns error when the remaining length exceeds src
    # packet
    -> mqtt.decode_fixed_header
