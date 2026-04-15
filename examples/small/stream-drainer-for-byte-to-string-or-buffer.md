# Requirement: "a library for reading the full contents of a byte stream into a string or a byte buffer"

Drains a stream handle by repeated reads into a growable buffer, then returns the accumulated bytes or the same bytes as a string.

std
  std.io
    std.io.read_chunk
      fn (stream: stream_handle, max: i32) -> result[bytes, string]
      + reads up to max bytes from the stream
      + returns empty bytes at end of stream
      - returns error on read failure
      # io

stream_drain
  stream_drain.read_all_bytes
    fn (stream: stream_handle) -> result[bytes, string]
    + reads the stream until end-of-stream and returns the concatenated bytes
    - returns error on any read failure
    ? callers are responsible for closing the stream
    # draining
    -> std.io.read_chunk
  stream_drain.read_all_string
    fn (stream: stream_handle) -> result[string, string]
    + reads the stream until end-of-stream and returns the bytes as a string
    - returns error on any read failure
    # draining
    -> std.io.read_chunk
