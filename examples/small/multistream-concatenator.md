# Requirement: "a library that concatenates multiple readable streams into a single readable stream"

Reads each source fully before advancing to the next, presenting one continuous stream.

std: (all units exist)

multistream
  multistream.new
    fn (sources: list[bus_state]) -> bus_state
    + constructs a stream that reads sources in order
    + accepts an empty source list and yields end-of-stream immediately
    # construction
  multistream.read
    fn (stream: bus_state, max: i32) -> result[bytes, string]
    + returns up to max bytes, advancing to the next source on end-of-stream
    - returns error from the current source without advancing
    # io
  multistream.close
    fn (stream: bus_state) -> result[void, string]
    + closes any remaining sources and surfaces the first error
    # lifecycle
