# Requirement: "a library that combines a writable stream and a readable stream into a single duplex stream"

A thin adapter: writes are forwarded to the writable side, reads are pulled from the readable side, and end-of-stream on either half closes the duplex.

std: (all units exist)

duplexify
  duplexify.wrap
    @ (writable: writable_stream, readable: readable_stream) -> duplex_stream
    + returns a duplex whose writes target writable and whose reads pull from readable
    # construction
  duplexify.write
    @ (d: duplex_stream, chunk: bytes) -> result[duplex_stream, string]
    + forwards chunk to the writable side
    - returns error when the duplex is already ended
    # write
  duplexify.read
    @ (d: duplex_stream, max_bytes: i32) -> result[tuple[bytes, duplex_stream], string]
    + returns up to max_bytes pulled from the readable side
    + returns empty bytes when the readable side is exhausted
    # read
  duplexify.end
    @ (d: duplex_stream) -> duplex_stream
    + marks the writable side as finished
    # lifecycle
  duplexify.destroy
    @ (d: duplex_stream, reason: string) -> duplex_stream
    + tears down both halves and records the reason
    # lifecycle
