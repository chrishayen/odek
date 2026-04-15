# Requirement: "a library that combines an array of streams into a single duplex stream"

Pipes writes into the first member and reads from the last, propagating close and errors.

std: (all units exist)

pipeline
  pipeline.new
    fn (stages: list[bus_state]) -> result[bus_state, string]
    + wires stages end-to-end into one duplex stream
    - returns error when stages is empty
    # construction
  pipeline.write
    fn (pipe: bus_state, chunk: bytes) -> result[i32, string]
    + forwards the chunk to the head stage and returns bytes accepted
    - returns error when any stage has failed
    # io
  pipeline.read
    fn (pipe: bus_state, max: i32) -> result[bytes, string]
    + returns up to max bytes from the tail stage
    - returns error when any stage has failed
    # io
  pipeline.close
    fn (pipe: bus_state) -> result[void, string]
    + closes every stage in order and surfaces the first error
    # lifecycle
