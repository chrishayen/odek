# Requirement: "a library for wrapping an in-memory value as a readable byte stream"

Turns various in-memory sources into an opaque stream handle that yields bytes on demand.

std: (all units exist)

to_stream
  to_stream.from_bytes
    fn (data: bytes) -> stream_handle
    + returns a stream that yields the bytes in order
    # construction
  to_stream.from_string
    fn (text: string) -> stream_handle
    + returns a stream over the UTF-8 encoding of text
    # construction
  to_stream.from_chunks
    fn (chunks: list[bytes]) -> stream_handle
    + returns a stream that yields each chunk in order with no boundary framing
    # construction
  to_stream.read
    fn (stream: stream_handle, max: i32) -> tuple[bytes, stream_handle]
    + returns up to max bytes from the stream along with the advanced handle
    + returns an empty slice when the stream is exhausted
    # consumption
  to_stream.drain
    fn (stream: stream_handle) -> bytes
    + returns all remaining bytes concatenated
    # consumption
