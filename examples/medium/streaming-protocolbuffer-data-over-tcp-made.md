# Requirement: "a streaming protocol buffer transport over TCP"

Length-prefixes each protobuf message on the wire so readers can frame them out of a byte stream. The library handles framing and dispatch; the socket and protobuf codec are pluggable.

std: (all units exist)

pbstream
  pbstream.encode_frame
    @ (payload: bytes) -> bytes
    + prefixes the payload with its length as a varint
    # framing
  pbstream.decode_frame
    @ (state: frame_state, chunk: bytes) -> result[tuple[frame_state, list[bytes]], string]
    + appends the chunk and returns any completed frames
    - returns error when a declared length is larger than the configured maximum
    ? partial frames remain buffered in state
    # framing
  pbstream.new_reader
    @ (conn: tcp_conn, max_frame_size: i32) -> reader_state
    + constructs a reader bound to a connection
    # construction
  pbstream.read_message
    @ (state: reader_state) -> result[tuple[reader_state, bytes], string]
    + returns the next complete message payload
    - returns error when the connection closes mid-frame
    # reader
  pbstream.new_writer
    @ (conn: tcp_conn) -> writer_state
    + constructs a writer bound to a connection
    # construction
  pbstream.write_message
    @ (state: writer_state, payload: bytes) -> result[writer_state, string]
    + frames the payload and writes it to the connection
    - returns error when the connection fails mid-write
    # writer
  pbstream.serve
    @ (listener: tcp_listener, handler: message_handler, max_frame_size: i32) -> result[void, string]
    + accepts connections and invokes the handler for each framed message
    - returns error when the listener fails
    # server
