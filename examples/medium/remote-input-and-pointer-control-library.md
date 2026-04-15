# Requirement: "a library for remotely controlling a host's pointer and keyboard from a client device"

A protocol layer: clients send input intents, the host applies them to an input device handle. Transport is abstracted.

std: (all units exist)

remote_input
  remote_input.encode_intent
    fn (intent: input_intent) -> bytes
    + serializes a pointer move, click, scroll, key press, or text entry intent as a compact frame
    # protocol
  remote_input.decode_intent
    fn (frame: bytes) -> result[input_intent, string]
    + parses a frame back into a typed intent
    - returns error on truncated frames
    - returns error on unknown intent tags
    # protocol
  remote_input.new_session
    fn (device: device_handle, pointer_speed: f32) -> session_state
    + creates a session bound to a device handle with a pointer speed multiplier
    # construction
  remote_input.apply_intent
    fn (session: session_state, intent: input_intent) -> result[session_state, string]
    + applies the intent to the underlying device and returns the updated session
    - returns error when the device refuses the operation
    # execution
  remote_input.type_text
    fn (session: session_state, text: string) -> result[session_state, string]
    + expands a text string to a sequence of key press/release intents and applies them
    # execution
